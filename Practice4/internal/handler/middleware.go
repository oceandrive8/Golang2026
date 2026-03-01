package handler

import (
	"log"
	"net/http"
	"time"
)

const validAPIKey = "Oceandrive"

// ------------------ LOGGING MIDDLEWARE ------------------

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		timestamp := time.Now().Format(time.RFC3339)
		log.Printf("[%s] %s %s", timestamp, r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// ------------------ AUTHENTICATION MIDDLEWARE ------------------

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey != validAPIKey {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
