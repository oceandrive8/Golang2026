package service

import (
	"errors"
	"testing"

	"Practice8/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{
		ID:   1,
		Name: "Bakytzhan Agai",
	}

	mockRepo.EXPECT().
		GetUserByID(1).
		Return(user, nil)

	result, err := userService.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{
		ID:   1,
		Name: "Bakytzhan Agai",
	}

	mockRepo.EXPECT().
		CreateUser(user).
		Return(nil)

	err := userService.CreateUser(user)

	assert.NoError(t, err)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Ayala"}

	mockRepo.EXPECT().
		GetByEmail("ayala@mail.com").
		Return(user, nil)

	err := service.RegisterUser(user, "ayala@mail.com")

	assert.Error(t, err)
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Ayala"}

	mockRepo.EXPECT().
		GetByEmail("ayala@mail.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		CreateUser(user).
		Return(nil)

	err := service.RegisterUser(user, "ayala@mail.com")

	assert.NoError(t, err)
}

func TestRegisterUser_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Ayala"}

	mockRepo.EXPECT().
		GetByEmail("ayala@mail.com").
		Return(nil, nil)

	mockRepo.EXPECT().
		CreateUser(user).
		Return(errors.New("db error"))

	err := service.RegisterUser(user, "ayala@mail.com")

	assert.Error(t, err)
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	err := service.UpdateUserName(1, "")

	assert.Error(t, err)
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	mockRepo.EXPECT().
		GetUserByID(1).
		Return(nil, errors.New("not found"))

	err := service.UpdateUserName(1, "Bota")

	assert.Error(t, err)
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Ayala"}

	mockRepo.EXPECT().
		GetUserByID(1).
		Return(user, nil)

	mockRepo.EXPECT().
		UpdateUser(user).
		Return(nil)

	err := service.UpdateUserName(1, "Bota")

	assert.NoError(t, err)
	assert.Equal(t, "Bota", user.Name)
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Ayala"}

	mockRepo.EXPECT().
		GetUserByID(1).
		Return(user, nil)

	mockRepo.EXPECT().
		UpdateUser(user).
		Return(errors.New("db error"))

	err := service.UpdateUserName(1, "Bota")

	assert.Error(t, err)
}

func TestDeleteUser_AdminBlocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	err := service.DeleteUser(1)

	assert.Error(t, err)
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	mockRepo.EXPECT().
		DeleteUser(2).
		Return(nil)

	err := service.DeleteUser(2)

	assert.NoError(t, err)
}

func TestDeleteUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	mockRepo.EXPECT().
		DeleteUser(2).
		Return(errors.New("db error"))

	err := service.DeleteUser(2)

	assert.Error(t, err)
}
