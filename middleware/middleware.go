package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/peteradeojo/lamp-logger/handlers"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		log.Println(r.Method, r.URL.Path, time.Since(start), wrapped.statusCode)
	})
}

func generateHMAC(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateHmac(next http.Handler) http.Handler {
	key := os.Getenv("ADMIN_HOST_KEY")

	prod := os.Getenv("GOENV") == "production"

	if key == "" {
		panic(errors.New("no secret key  connfigured in the application"))
	}

	secretKey := []byte(key)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if prod {
			signature := r.Header.Get("X-Signature")

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "unable to read body", http.StatusInternalServerError)
				return
			}

			if signature == "" {
				handlers.RespondWithError(w, 401, "no signature provided")
				return
			}

			// Log or process the body (example: printing it)
			// fmt.Printf("Request Body: %s\n", string(body))

			hash := hmac.New(sha256.New, secretKey)
			_, err = hash.Write(body)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			expected := hex.EncodeToString(hash.Sum(nil))

			check := hmac.Equal([]byte(signature), []byte(expected))

			if !check {
				handlers.RespondWithError(w, http.StatusForbidden, "Forbidden")
				return
			}

			// Restore the body for the next handler.
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(w, r)
	})
}
