package usecase

import (
	"awesomeProject/internal/repository/_postgres/users"
	"awesomeProject/pkg/modules"
)

type UserUsecase struct {
	repo *users.Repository
}

func NewUserUsecase(repo *users.Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Get all users
func (u *UserUsecase) GetAllUsers(limit, offset int) ([]modules.User, error) {
	return u.repo.GetUsers(limit, offset)
}

// Get user by ID
func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

// Create new user
func (u *UserUsecase) CreateUser(user *modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

// Update user
func (u *UserUsecase) UpdateUser(user *modules.User) error {
	return u.repo.UpdateUser(user)
}

// Delete user
func (u *UserUsecase) DeleteUser(id int) (int64, error) {
	return u.repo.DeleteUser(id)
}
