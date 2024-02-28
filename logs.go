package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/peteradeojo/lamp-logger/handlers"
	"github.com/peteradeojo/lamp-logger/internal/database"
)

func (apiCfg *ApiConfig) saveLog(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Text    string      `json:"text"`
		Level   string      `json:"level"`
		Context interface{} `json:"context"`
		Tags    interface{} `json:"tags"`
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

	ip := r.RemoteAddr
	ipAddress := sql.NullString{
		String: ip,
		Valid:  false,
	}
	if ip == "" {
		ipAddress.String = ip
		ipAddress.Valid = true
	}

	tags, err := json.Marshal(params.Tags)
	if err != nil {
		log.Println(err)
		handlers.RespondWithError(w, 400, "Unable to parse tags")
		return
	}

	_, err = apiCfg.DB.SaveLogs(r.Context(), database.SaveLogsParams{
		Apptoken: appId,
		Text:     params.Text,
		Level:    params.Level,
		Context:  context,
		Ip:       ipAddress,
		Tags:     tags,
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

// func (apiCfg *ApiConfig) getLogs(w http.ResponseWriter, r *http.Request) {
// 	token := chi.URLParam(r, "token")
// 	if token == "" {
// 		handlers.RespondWithError(w, 404, "")
// 		return
// 	}

// 	logs, err := apiCfg.DB.GetLogs(r.Context(), token)
// 	if err != nil {
// 		handlers.RespondWithError(w, 500, "An error occured: "+err.Error())
// 		return
// 	}

// 	handlers.Respond(w, 200, logs)
// }
