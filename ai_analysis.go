package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/peteradeojo/lamp-logger/handlers"

	llog "log"
)

type ModelResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

type GenerateResponse struct {
	ModelResponse
	Response string `json:"response"`
}

type ChatResponse struct {
	ModelResponse
	Message ChatMessage `json:"message"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (apiCfg *ApiConfig) RunAIAnalysis(w http.ResponseWriter, r *http.Request) {
	logId, err := strconv.Atoi(chi.URLParam(r, "logId"))
	if err != nil {
		handlers.SaveError(apiCfg.DB, r, err)
		handlers.RespondWithError(w, 400, "Malformed request")
		return
	}

	log, err := apiCfg.DB.GetLog(r.Context(), int64(logId))
	if err != nil {
		llog.Println(err)
		handlers.RespondWithError(w, 500, "Unable to fetch log")
		return
	}

	content := fmt.Sprintf("**Message**: %s\n**Context**: %s", log.Text, log.Context.RawMessage)

	var reader = strings.NewReader(fmt.Sprintf(`{"messages": [{"role": "system", "content": "You are to act as an analyser for log messages. You will be provided with a generated log message from a running application. The message might contain context provided as a JSON string. You will attempt to determine the programming language from the log message and provided context and attempt to provide a description about the issue to be diagnosed and provide pointers on how best to resolve the issue."}, {"role": "user", "content": "%v"}], "stream": false, "model": "llama3"}`, content))

	resp, err := http.Post("http://127.0.0.1:11434/api/chat", "application/json", reader)

	if err != nil {
		llog.Println(err)
		handlers.SaveError(apiCfg.DB, r, err)
		handlers.RespondWithError(w, 500, err.Error())
		return
	}
	defer resp.Body.Close()

	i, err := io.ReadAll(resp.Body)
	if err != nil {
		handlers.SaveError(apiCfg.DB, r, err)
		llog.Println(err)
		handlers.RespondWithError(w, 500, err.Error())
		return
	}

	var data ChatResponse
	err = json.Unmarshal(i, &data)
	if err != nil {
		handlers.SaveError(apiCfg.DB, r, err)
		llog.Println(err)
		handlers.RespondWithError(w, 500, err.Error())
		return
	}

	handlers.Respond(w, 200, data)
}
