package main

import (
	"context"
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
}

func (a *ApiConfig) ProcessEmailQueue() {
	ctx := context.Background()

	l, err := a.redisClient.LLen(ctx, "notifications_queue").Result()
	if err != nil {
		log.Println(err)
	}

	var mail MailMesage
	for i := 0; i < int(l); i += 1 {
		record, err := a.redisClient.LPop(ctx, "notifications_queue").Result()
		if err != nil {
			log.Println(err)
			continue
		}

		err = json.Unmarshal([]byte(record), &mail)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("Mail: %v", mail)

		err = SendMail(mail)
		if err != nil {
			log.Println(err)
			a.redisClient.RPush(ctx, "notifications_queue", record)
			continue
		}
	}
}
