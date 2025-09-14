package dto

import "coffee_shop/models"

type RegisterResponse struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
}

func ToRegisterResponse(user models.User) RegisterResponse {
	return RegisterResponse{
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
	}
}

type LoginResponse struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	Token       string `json:"token"`
}

func ToLoginResponse(user models.User, token string) LoginResponse {
	return LoginResponse{
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
		Token:       token,
	}
}
