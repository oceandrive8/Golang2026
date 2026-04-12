package usecase

import (
	"Practice7/internal/entity"
	"Practice7/internal/usecase/repo"
	"Practice7/utils"
	"fmt"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {

	userFromRepo, err := u.repo.LoginUser(input)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !utils.CheckPasswordHash(userFromRepo.Password, input.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {

	createdUser, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}

	sessionID := uuid.New().String()

	return createdUser, sessionID, nil
}
func (u *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UserUseCase) PromoteUser(id string) error {
	return u.repo.PromoteUser(id)
}
