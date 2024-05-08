package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/peteradeojo/lamp-logger/internal/database"
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

func SaveError(db *database.Queries, r *http.Request, err error) {
	if err != nil {
		db.CreateSystemLog(r.Context(), database.CreateSystemLogParams{
			Text:  err.Error(),
			Level: database.LogLevelError,
			Stack: sql.NullString{
				Valid:  true,
				String: "Unavailable",
			},
		})
	}
}
