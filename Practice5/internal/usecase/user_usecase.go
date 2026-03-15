package usecase

import (
	"awesomeProject/internal/repository/_postgres/users"
	"awesomeProject/pkg/modules"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserUsecase struct {
	repo *users.Repository
}

func NewUserUsecase(repo *users.Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Get all users
func (u *UserUsecase) GetAllUsers(limit, offset int, orderBy []string) ([]modules.User, error) {
	return u.repo.GetUsers(limit, offset, orderBy)
}
func (u *UserUsecase) GetPaginatedUsers(
	page int,
	pageSize int,
	filters map[string]interface{},
	orderBy []string,
) (modules.PaginatedResponse, error) {
	return u.repo.GetPaginatedUsers(page, pageSize, filters, orderBy)
}

// Get user by ID
func (u *UserUsecase) GetUserByID(id uuid.UUID) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

// Create new user
func (u *UserUsecase) CreateUser(user *modules.User) (uuid.UUID, error) {
	return u.repo.CreateUser(user)
}

// Update user
func (u *UserUsecase) UpdateUser(user *modules.User) error {
	return u.repo.UpdateUser(user)
}

// Delete user
func (u *UserUsecase) DeleteUser(id uuid.UUID) (int64, error) {
	return u.repo.DeleteUser(id)
}
func (u *UserUsecase) GetCommonFriends(ctx context.Context, user1, user2 uuid.UUID) ([]modules.User, error) {
	friends, err := u.repo.GetCommonFriends(ctx, user1, user2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch common friends: %w", err)
	}
	// Do NOT return an error if friends is empty
	return friends, nil
}
