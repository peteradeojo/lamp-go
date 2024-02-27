package main

import (
	"encoding/json"
	"net/http"

	"github.com/peteradeojo/lamp-logger/handlers"
	"github.com/peteradeojo/lamp-logger/internal/database"
)

func (apiCfg *ApiConfig) saveLog(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Text    string      `json:"text"`
		Level   string      `json:"level"`
		Context interface{} `json:"context"`
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

	context, err := json.Marshal(params.Context)
	if err != nil {
		handlers.RespondWithError(w, 400, err.Error())
		return
	}

	_, err = apiCfg.DB.SaveLogs(r.Context(), database.SaveLogsParams{
		Apptoken: appId,
		Text:     params.Text,
		Level:    params.Level,
		Context:  context,
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
