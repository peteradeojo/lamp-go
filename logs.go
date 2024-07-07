package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/peteradeojo/lamp-logger/handlers"
	"github.com/peteradeojo/lamp-logger/internal/database"
	"github.com/sqlc-dev/pqtype"
)

type Log struct {
	ID        int                   `json:"id"`
	Text      string                `json:"text"`
	Apptoken  string                `json:"apptoken"`
	Level     string                `json:"level"`
	Createdat sql.NullTime          `json:"createdat"`
	Context   pqtype.NullRawMessage `json:"context"`
	Ip        sql.NullString        `json:"ip"`
	Tags      pqtype.NullRawMessage `json:"tags"`
}

type ExportJob struct {
	Status   string `json:"status"`
	Path     string `json:"path"`
	AppToken string `json:"token"`
}

func (apiCfg *ApiConfig) saveLog(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Text    string        `json:"text"`
		Level   string        `json:"level"`
		Context []interface{} `json:"context"`
		Tags    []string      `json:"tags"`
		Ip      string        `json:"ip"`
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

	// context, err := json.Marshal(params.Context)
	// if err != nil {
	// 	handlers.RespondWithError(w, 400, err.Error())
	// 	return
	// }

	ip := r.RemoteAddr
	// ipAddress := sql.NullString{
	// 	String: ip,
	// 	Valid:  true,
	// }

	params.Ip = ip

	// tags, err := json.Marshal(params.Tags)
	// if err != nil
	// 	log.Println(err)
	// 	handlers.RespondWithError(w, 400, "Unable to parse tags")
	// 	return
	// }

	lLog, err := json.Marshal(struct {
		AppToken  string        `json:"apptoken"`
		Tags      []string      `json:"tags"`
		Context   []interface{} `json:"context"`
		CreatedAt string        `json:"createdat"`
		Level     string        `json:"level"`
		Text      string        `json:"text"`
		Ip        string        `json:"ip"`
	}{
		AppToken:  appId,
		Tags:      params.Tags,
		Context:   params.Context,
		CreatedAt: time.Now().Format(time.DateTime),
		Level:     params.Level,
		Text:      params.Text,
		Ip:        params.Ip,
	})

	if err != nil {
		handlers.RespondWithError(w, 400, err.Error())
		return
	}

	// _, err = apiCfg.DB.SaveLogs(r.Context(), database.SaveLogsParams{
	// 	Apptoken: appId,
	// 	Text:     params.Text,
	// 	Level:    params.Level,
	// 	Context:  pqtype.NullRawMessage{RawMessage: context, Valid: true},
	// 	Ip:       ipAddress,
	// 	// Tags:     tags,
	// 	Tags: pqtype.NullRawMessage{RawMessage: tags, Valid: true},
	// })
	i := apiCfg.redisClient.RPush(r.Context(), "pending_logs", lLog)
	if i.Err() != nil {
		apiCfg.DB.CreateSystemLog(r.Context(), database.CreateSystemLogParams{
			Text:  i.Err().Error(),
			Level: database.LogLevelError,
			Stack: sql.NullString{
				Valid:  true,
				String: "Unavailable",
			},
		})
		handlers.RespondWithError(w, 500, i.Err().Error())
		return
	}

	err = apiCfg.sendSocketMessage("log", appId, string(lLog))
	if err != nil {
		reportError(r.Context(), err, pqtype.NullRawMessage{
			Valid:      false,
			RawMessage: nil,
		})
	}

	handlers.Respond(w, 200, struct {
		Message string `json:"message"`
	}{
		Message: "Log saved successfully",
	})
}

func (apiCfg *ApiConfig) exportLogs(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		App string `json:"app"`
	}

	params := &parameters{}

	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handlers.RespondWithError(w, 500, err.Error())
		return
	}

	date := time.Now().Format("2006-01-02")

	go apiCfg.generateLogExport(context.Background(), params.App, fmt.Sprintf("exports/%s/%s/Book1.xlsx", date, params.App))

	handlers.Respond(w, 200, handlers.ApiResponse{Message: "Exporting generated file."})
}
