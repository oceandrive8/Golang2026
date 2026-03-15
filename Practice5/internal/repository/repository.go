package repository

import (
	"awesomeProject/internal/repository/_postgres"
	"awesomeProject/internal/repository/_postgres/users"
	"awesomeProject/pkg/modules"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUsers(limit, offset int, orderBy []string) ([]modules.User, error)
	GetUserByID(id uuid.UUID) (*modules.User, error)
	CreateUser(user *modules.User) (uuid.UUID, error)
	UpdateUser(user *modules.User) error
	DeleteUser(id uuid.UUID) (int64, error)
	GetPaginatedUsers(page int, pageSize int, filters map[string]interface{}, orderBy []string) (modules.PaginatedResponse, error)
	GetCommonFriends(ctx context.Context, user1, user2 uuid.UUID) ([]modules.User, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
