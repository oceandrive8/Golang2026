package app

import (
	"awesomeProject/internal/handler"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/repository/_postgres"
	"awesomeProject/internal/repository/_postgres/users"
	"awesomeProject/internal/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	userRepo, ok := repos.UserRepository.(*users.Repository)
	if !ok {
		return fmt.Errorf("invalid user repository type")
	}

	// ------------------ USECASE ------------------
	userUsecase := usecase.NewUserUsecase(userRepo)

	// ------------------ HANDLER ------------------
	userHandler := handler.NewUserHandler(userUsecase)

	// ------------------ HTTP SERVER ------------------
	mux := http.NewServeMux()

	// ------------------ Routes ------------------

	// GET all users (limit/offset, optional sorting)
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetAllUsers(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// GET paginated users: /users/paginated?page=1&pageSize=10
	mux.HandleFunc("/users/paginated", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetUsersPaginated(w, r)
		} else if r.Method == http.MethodPost {
			userHandler.CreateUser(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// User item routes: GET by ID, PATCH/PUT, DELETE
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

	// Common friends
	mux.HandleFunc("/users/common_friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetCommonFriends(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// ------------------ Middleware ------------------
	finalHandler := handler.LoggingMiddleware(handler.AuthMiddleware(mux))

	log.Println("Server running on :8080")
	return http.ListenAndServe(":8080", finalHandler)
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		Username:    getEnv("DB_USER", "ayalanurakyn"),
		Password:    getEnv("DB_PASSWORD", "12345678"),
		DBName:      getEnv("DB_NAME", "mydb"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
