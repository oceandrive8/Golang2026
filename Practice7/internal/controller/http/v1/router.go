package v1

import (
	"Practice7/internal/usecase"
	"Practice7/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.RouterGroup, u usecase.UserInterface) {
	log := logger.New()

	newUserRoutes(handler, u, log)
}
