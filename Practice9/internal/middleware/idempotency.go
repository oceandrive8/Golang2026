package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"Practice9/internal/storage"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Idempotency(store *storage.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		rec, exists := store.GetOrCreate(key)

		if exists {
			if rec.Status == "processing" {
				fmt.Println("Conflict for key:", key)
				http.Error(w, "conflict", http.StatusConflict)
				return
			}

			fmt.Println("Returning cached response for key:", key)
			w.WriteHeader(rec.StatusCode)
			w.Write(rec.Body)
			return
		}

		fmt.Println("Processing started for key:", key)

		recorder := &responseRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		store.SetCompleted(key, recorder.status, recorder.body.Bytes())
	})
}
