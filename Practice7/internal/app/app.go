package app

import (
	"Practice7/internal/controller/http/v1"
	"Practice7/internal/usecase"
	"Practice7/internal/usecase/repo"
	"Practice7/pkg/postgres"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=auth_service"

	pg, err := postgres.NewPostgres(dsn)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	userRepo := repo.NewUserRepo(pg)
	userUsecase := usecase.NewUserUseCase(userRepo)

	router := gin.Default()

	api := router.Group("/v1")
	{
		v1.NewRouter(api, userUsecase)
	}

	router.Run(":8090")
}
