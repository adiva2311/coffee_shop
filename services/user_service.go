package services

import (
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/repositories"
	"coffee_shop/utils"
	"errors"

	"gorm.io/gorm"
)

type UserService interface {
	// Define user-related business logic methods here
	Register(registerRequest models.User) (dto.RegisterResponse, error)
	Login(email, password string) (dto.LoginResponse, error)
	UpdateUser(user_id uint, user models.User) (models.User, error)
	DeleteUser(user_id uint) (int, error)
	GetUserByID(user_id uint) ([]models.User, error)
}

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
func (u *UserServiceImpl) Login(email string, password string) (dto.LoginResponse, error) {
	// // Check if the email exists
	// user, err := u.userRepo.CheckEmailValid(email)
	// if err != nil {
	// 	return dto.LoginResponse{}, errors.New("invalid email or password")
	// }

	// // Check if the password is correct
	// if !utils.CheckPasswordHash(password, user.Password) {
	// 	return dto.LoginResponse{}, errors.New("invalid email or password")
	// }

	panic("unimplemented")
}

// DeleteUser implements UserService.
func (u *UserServiceImpl) DeleteUser(user_id uint) (int, error) {
	panic("unimplemented")
}

// GetUserByID implements UserService.
func (u *UserServiceImpl) GetUserByID(user_id uint) ([]models.User, error) {
	panic("unimplemented")
}

// UpdateUser implements UserService.
func (u *UserServiceImpl) UpdateUser(user_id uint, user models.User) (models.User, error) {
	panic("unimplemented")
}

func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		userRepo: repositories.NewUserRepository(db),
	}
}
