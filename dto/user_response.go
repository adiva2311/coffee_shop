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
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	PhoneNumber    string `json:"phone_number"`
	AccessToken    string `json:"access_token"`
	RefresherToken string `json:"refresher_token"`
}

func ToLoginResponse(user models.User, accessToken string, refresherToken string) LoginResponse {
	return LoginResponse{
		Name:           user.Name,
		Email:          user.Email,
		Role:           user.Role,
		PhoneNumber:    user.PhoneNumber,
		AccessToken:    accessToken,
		RefresherToken: refresherToken,
	}
}
