package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"awesomeProject/internal/handler"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/repository/_postgres"
	"awesomeProject/internal/repository/_postgres/users"
	"awesomeProject/internal/usecase"
	"awesomeProject/pkg/modules"
)

func Run() error {
	// ------------------ DATABASE ------------------
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbConfig := initPostgreConfig()
	fmt.Printf("Connecting to DB: %+v\n", dbConfig)

	pg, err := _postgres.NewPGXDialectSafe(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}
	fmt.Println("Connected to database successfully")

	// Run migrations
	if err := _postgres.AutoMigrate(dbConfig); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	fmt.Println("Migrations completed successfully")

	// ------------------ REPOSITORY ------------------
	repos := repository.NewRepositories(pg)
	userRepo := repos.UserRepository.(*users.Repository)

	// ------------------ USECASE ------------------
	userUsecase := usecase.NewUserUsecase(userRepo)

	// ------------------ HANDLER ------------------
	userHandler := handler.NewUserHandler(userUsecase)

	// ------------------ HTTP SERVER ------------------
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetAllUsers(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByID(w, r)
		case http.MethodPatch, http.MethodPut:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Wrap mux with middleware (Auth first, then Logging)
	finalHandler := handler.LoggingMiddleware(handler.AuthMiddleware(mux))

	log.Println("Server running on :8080")
	return http.ListenAndServe(":8080", finalHandler)
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "ayalanurakyn",
		Password:    "",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
