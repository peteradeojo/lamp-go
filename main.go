package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/peteradeojo/lamp-logger/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("No PORT")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("No DB_URL")
	}

	dbCxn, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal("Unable to open database connection: ", err)
	}

	defer dbCxn.Close()

	apiCfg := ApiConfig{
		DB: database.New(dbCxn),
	}

	router := chi.NewRouter()

	v1Router := chi.NewRouter()
	v1Router.Post("/logs", apiCfg.saveLog)

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
