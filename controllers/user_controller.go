package controllers

import (
	"coffee_shop/config"
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/services"
	"coffee_shop/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController interface {
	// Define user-related controller methods here
	Register(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	RefreshToken(c echo.Context) error
}

type UserControllerImpl struct {
	UserService services.UserService
}

var ctx = context.Background()

// Define user-related controller methods here
func (u *UserControllerImpl) Register(c echo.Context) error {
	userPayload := new(dto.RegisterRequest)
	err := c.Bind(userPayload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload : " + err.Error(),
		})
	}

	// Capitalize the first letter of the Name
	nameStr := strings.ToTitle(userPayload.Name)

	if userPayload.PhoneNumber == "" {
		userPayload.PhoneNumber = "XXXXXX"
	}

	result, err := u.UserService.Register(models.User{
		Name:        nameStr,
		Email:       userPayload.Email,
		Password:    userPayload.Password,
		Role:        userPayload.Role,
		PhoneNumber: userPayload.PhoneNumber,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to register user : " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User registered successfully",
		Data:    result,
	}

	return c.JSON(http.StatusOK, apiResponse)
}

func (u *UserControllerImpl) Login(c echo.Context) error {
	userPayload := new(dto.LoginRequest)
	err := c.Bind(userPayload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload : " + err.Error(),
		})
	}

	result, err := u.UserService.Login(dto.LoginRequest{
		Email:    userPayload.Email,
		Password: userPayload.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  echo.ErrUnauthorized.Code,
			Message: "Failed to Login : " + err.Error() + " | " + echo.ErrUnauthorized.Error(),
		})
	}
	fmt.Println(result)

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User logged in successfully",
		Data:    result,
	}
	return c.JSON(http.StatusOK, apiResponse)
}

func (u *UserControllerImpl) Logout(c echo.Context) error {
	userEmail := c.Get("email").(string)
	// Invalidate the token by adding it to a blacklist in Redis
	log.Println("User Email from token:", userEmail)
	rdb, err := config.RedisClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to connect to Redis: " + err.Error(),
		})
	}
	err = rdb.Del(ctx, userEmail).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to logout: " + err.Error(),
		})
	}
	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User logged out successfully",
	}
	return c.JSON(http.StatusOK, apiResponse)
}

func (u *UserControllerImpl) UpdateUser(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *UserControllerImpl) DeleteUser(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *UserControllerImpl) GetUserByID(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *UserControllerImpl) RefreshToken(c echo.Context) error {
	req := new(dto.RefreshTokenRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload : " + err.Error(),
		})
	}

	// Parse refresh token
	// token, err := jwt.ParseWithClaims(req.RefresherToken, &utils.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return utils.GetSecretKey(), nil
	// })
	token, err := jwt.Parse(req.RefresherToken, func(token *jwt.Token) (interface{}, error) {
		return utils.GetSecretKey(), nil
	})
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid refresh token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, _ := claims["user_id"].(float64)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	newAccessToken, _ := utils.GenerateJWT(uint(userID), email, role)

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Token refreshed successfully",
		Data:    newAccessToken,
	}
	return c.JSON(http.StatusOK, apiResponse)
}

func NewUserController(db *gorm.DB) UserControllerImpl {
	service := services.NewUserService(db)
	return UserControllerImpl{
		UserService: service,
	}
}
