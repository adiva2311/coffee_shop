package services

import (
	"coffee_shop/config"
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/repositories"
	"coffee_shop/utils"
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	// Define user-related business logic methods here
	Register(registerRequest models.User) (dto.RegisterResponse, error)
	Login(request dto.LoginRequest) (dto.LoginResponse, error)
	UpdateUser(user_id uint, request models.User) (dto.UpdateUserResponse, error)
	DeleteUser(user_id uint) (int, error)
	GetUserByID(user_id uint) (dto.GetUserByIDResponse, error)
}

var ctx = context.Background()

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

// Register implements UserService.
func (u *UserServiceImpl) Register(registerRequest models.User) (dto.RegisterResponse, error) {
	// Check if the email already exists
	emailExist, _ := u.userRepo.CheckEmailExists(registerRequest.Email)
	if emailExist {
		return dto.RegisterResponse{}, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:        registerRequest.Name,
		Email:       registerRequest.Email,
		Password:    hashedPassword,
		Role:        registerRequest.Role,
		PhoneNumber: registerRequest.PhoneNumber,
	}

	err = u.userRepo.RegisterUser(user)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("failed to register user")
	}
	return dto.ToRegisterResponse(*user), nil
}

// Login implements UserService.
func (u *UserServiceImpl) Login(request dto.LoginRequest) (dto.LoginResponse, error) {
	// Check if the email exists
	user, err := u.userRepo.CheckEmailValid(request.Email)
	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid email")
	}

	// Check if the password is correct
	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return dto.LoginResponse{}, errors.New("invalid password")
	}

	// Generate JWT token
	accessToken, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refresherToken, err := utils.GenerateRefresherJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// Save Refresher Token to Redis
	rdb, err := config.RedisClient()
	if err != nil {
		log.Fatal("Failed Connect to Redis")
	}

	err = rdb.Set(ctx, user.Email, refresherToken, 7*24*time.Hour).Err()
	if err != nil {
		return dto.LoginResponse{}, errors.New("failed to save refresher token to redis")
	}

	return dto.ToLoginResponse(*user, accessToken, refresherToken), nil
}

// DeleteUser implements UserService.
func (u *UserServiceImpl) DeleteUser(user_id uint) (int, error) {
	rowsAffected, err := u.userRepo.DeleteUser(user_id)
	if err != nil {
		return 0, errors.New("failed to delete user: " + err.Error())
	}
	return rowsAffected, nil
}

// GetUserByID implements UserService.
func (u *UserServiceImpl) GetUserByID(user_id uint) (dto.GetUserByIDResponse, error) {
	userInfo, err := u.userRepo.GetUserByID(user_id)
	if err != nil {
		return dto.GetUserByIDResponse{}, errors.New("user not found : " + err.Error())
	}
	return dto.ToGetUserByIDResponse(*userInfo), nil
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(user_id uint, request models.User) (dto.UpdateUserResponse, error) {
	// Hash the password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return dto.UpdateUserResponse{}, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:        request.Name,
		Password:    hashedPassword,
		Role:        request.Role,
		PhoneNumber: request.PhoneNumber,
	}

	err = u.userRepo.UpdateUser(user_id, user)
	if err != nil {
		return dto.UpdateUserResponse{}, errors.New("failed to update user: " + err.Error())
	}

	return dto.ToUpdateUserResponse(*user), nil
}

func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		userRepo: repositories.NewUserRepository(db),
	}
}
