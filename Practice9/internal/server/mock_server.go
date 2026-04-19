package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func NewMockServer() *httptest.Server {
	counter := 0

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++

		if counter <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "temporary failure")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"success"}`)
	}))
}
