package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/peteradeojo/lamp-logger/handlers"
	"github.com/peteradeojo/lamp-logger/internal/database"
)

func (apiCfg *ApiConfig) saveLog(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Appid int64
		Text  string
		Level string
	}

	params := &parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		handlers.RespondWithError(w, 400, "Bad request")
		return
	}

	appId := r.Header.Get("APP_ID")
	if appId == "" {
		handlers.RespondWithError(w, 400, "Bad Request")
		return
	}

	app, err := apiCfg.DB.GetAppWithToken(r.Context(), sql.NullString{
		String: appId,
		Valid:  true,
	})

	if err != nil {
		handlers.RespondWithError(w, 400, "Invalid app token")
		return
	}

	_, err = apiCfg.DB.SaveLogs(r.Context(), database.SaveLogsParams{
		Appid: app.ID,
		Text:  params.Text,
		Level: params.Level,
	})
	if err != nil {
		handlers.RespondWithError(w, 500, err.Error())
		return
	}

	handlers.Respond(w, 200, struct {
		Message string `json:"message"`
	}{
		Message: "log saved successfully",
	})
}
