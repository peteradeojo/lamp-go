package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/peteradeojo/lamp-logger/handlers"
)

func (apiCfg *ApiConfig) SendNotification(w http.ResponseWriter, r *http.Request) {
	var data MailMesage

	err := handlers.ParseJson(r.Body, &data)

	if err != nil {
		log.Println(err)
		handlers.RespondWithError(w, 400, err.Error())

		return
	}

	if err := SendMail(data); err != nil {
		fmt.Println(err)
	}

	handlers.Respond(w, http.StatusOK, map[string]any{
		"success": true,
	})

	r.Body.Close()
}

func (apiCfg *ApiConfig) QueueEmail(w http.ResponseWriter, r *http.Request) {
	var data MailMesage
	err := handlers.ParseJson(r.Body, &data)

	if err != nil {
		log.Println(err)
		handlers.RespondWithError(w, 400, err.Error())

		return
	}

	j, err := json.Marshal(data)
	if err != nil {

		log.Println(err)
		handlers.RespondWithError(w, 500, err.Error())
	}

	apiCfg.redisClient.LPush(r.Context(), "notifications_queue", string(j))
	handlers.Respond(w, 200, map[string]string{
		"message": "Success",
	})

	return
}
