package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/peteradeojo/lamp-logger/internal/database"
	"github.com/sqlc-dev/pqtype"
)

type ApiConfig struct {
	DB          *database.Queries
	redisClient *redis.Client
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

	router := chi.NewRouter()

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

	router.Mount("/v1", v1Router)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	log.Printf("Server running on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func reportError(cx context.Context, err error, context pqtype.NullRawMessage) {
	apiCfg.DB.CreateSystemLog(cx, database.CreateSystemLogParams{
		Text:    fmt.Sprintf("Unable to register job: %v", err),
		Level:   "error",
		Context: context,
	})
}
