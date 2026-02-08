package main

import (
	"log"
	"net/http"

	"awesomeProject/internal/handlers"
	"awesomeProject/internal/middleware"
)

func main() {
	taskHandler := handlers.NewTaskHandler()

	var handler http.Handler = taskHandler
	handler = middleware.APIKey(handler)
	handler = middleware.Logging(handler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
