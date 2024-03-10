package main

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/peteradeojo/lamp-logger/internal/database"

	dotenv "github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
)

func TestMain(t *testing.T) {
	env := os.Getenv("ENVIRONMENT")
	if env != "production" {
		dotenv.Load(".env.test")
	}

	apiCfg := ApiConfig{}

	cxn, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_URL"))

	if err != nil {
		t.Fatal(err)
	}

	apiCfg.DB = database.New(cxn)
	apptoken := "6f7b8451-0724-492d-af75-a8da0e2b108f"

	test := apitest.New()
	test.HandlerFunc(apiCfg.saveLog)
	test.Request().Header("APP_ID", apptoken).Body(`{"level": "info", "text": "An error occured"}`)

	ex := test.Post("/logs").Expect(t)
	ex.Body(`{"message": "Log saved successfully"}`).Status(http.StatusOK)
	ex.End()
}
