package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/peteradeojo/lamp-logger/handlers"
)

func (apiCfg *ApiConfig) getApps(w http.ResponseWriter, r *http.Request) {
	apps, err := apiCfg.DB.GetApps(r.Context(), 20)
	if err != nil {
		handlers.RespondWithError(w, 500, err.Error())
		log.Println(err)
		return
	}

	handlers.Respond(w, 200, apps)
}

func (apiCfg *ApiConfig) getAppWithToken(w http.ResponseWriter, r *http.Request) {
	appToken := chi.URLParam(r, "app")
	token := sql.NullString{
		Valid: false,
	}
	if appToken == "" {
		handlers.RespondWithError(w, 400, "Bad request")
		return
	} else {
		token.String = appToken
		token.Valid = true
	}

	apps, err := apiCfg.DB.GetAppWithToken(r.Context(), token)
	if err != nil {
		handlers.RespondWithError(w, 500, err.Error())
		return
	}

	handlers.Respond(w, 200, apps)
}
