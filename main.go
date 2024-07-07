package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/peteradeojo/lamp-logger/internal/database"
	"github.com/sqlc-dev/pqtype"
	"github.com/zishang520/socket.io/v2/socket"
)

type ApiConfig struct {
	DB          *database.Queries
	redisClient *redis.Client
	ioClient    *socket.Server
}

var apiCfg *ApiConfig

func main() {
	apiCfg = &ApiConfig{
		// DB:          database.New(dbCxn),
		// redisClient: redisClient,
	}

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("No PORT")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("No DB_URL")
	}

	dbCxn, err := sql.Open(os.Getenv("DB_TYPE"), dbUrl)
	if err != nil {
		log.Fatal("Unable to open database connection: ", err)
	}

	apiCfg.DB = database.New(dbCxn)

	defer dbCxn.Close()

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	if redisClient == nil {
		log.Fatal("Unable to create redis client")
	}

	apiCfg.redisClient = redisClient
	defer redisClient.Close()

	// Websocket socketio setup
	router := chi.NewRouter()

	io := socket.NewServer(nil, nil)

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://*"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Post("/logs", apiCfg.saveLog)
	v1Router.Post("/export", apiCfg.exportLogs)
	v1Router.Get("/apps", apiCfg.getApps)
	v1Router.Get("/apps/{app}", apiCfg.getAppWithToken)

	apiCfg.BootstrapWs(io)

	router.Mount("/v1", v1Router)
	router.Mount("/socket.io", io.ServeHandler(nil))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	go func() {
		log.Printf("Server running on port %s", portString)
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	io.Close(nil)

	log.Println("server stopped")

}

func reportError(cx context.Context, err error, context pqtype.NullRawMessage) {
	apiCfg.DB.CreateSystemLog(cx, database.CreateSystemLogParams{
		Text:    err.Error(),
		Level:   "error",
		Context: context,
	})
}

func (apiCfg *ApiConfig) sendSocketMessage(evt, to string, message any) error {
	return apiCfg.ioClient.To(socket.Room(to)).Emit(evt, message)
}

func (api *ApiConfig) BootstrapWs(io *socket.Server) {
	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		client.On("connect-log-stream", func(a ...any) {
			token := a[0]
			if room, ok := token.(string); ok {
				client.Join(socket.Room(room))
				fmt.Println("Joined room")
			} else {
				fmt.Println("Token is not of type string")
			}
		})

		client.On("disconnect", func(a ...any) {
			fmt.Println("disconnected")
		})

		log.Printf("connected: %s\n", client.Id())
	})

	api.ioClient = io
}
