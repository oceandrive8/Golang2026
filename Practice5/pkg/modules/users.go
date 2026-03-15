package modules

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Email    *string    `json:"email"`
	Gender   string     `json:"gender"`
	Birthday *time.Time `json:"birthDate,omitempty"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}
