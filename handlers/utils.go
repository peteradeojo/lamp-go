package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func Respond(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Println("Unable to marshal json data from", payload)
		Respond(w, 500, struct{ message string }{message: err.Error()})
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Println(err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type errMessage struct {
		Message string `json:"message"`
	}

	if code > 499 {
		log.Printf("Responding with 500 level error: %v", msg)
	}

	message := errMessage{
		Message: msg,
	}
	Respond(w, code, message)
}
