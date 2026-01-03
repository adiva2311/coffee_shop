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
	"time"

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

	// Store the token in Redis with an expiration time
	rdb, err := config.RedisClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to connect to Redis: " + err.Error(),
		})
	}

	// Store Access Token and Refresh Token in Redis
	err = rdb.Set(ctx, userPayload.Email, result.AccessToken, time.Minute*15).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Store Access Token to Redis: " + err.Error(),
		})
	}

	err = rdb.Set(ctx, "refresh:"+userPayload.Email, result.RefresherToken, time.Hour*24).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to Store Refresh Token to Redis: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User logged in successfully",
		Data:    result,
	}
	return c.JSON(http.StatusOK, apiResponse)
}

func (u *UserControllerImpl) Logout(c echo.Context) error {
	userEmail, ok := c.Get("email").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	// Delete tokens from Redis
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

	err = rdb.Del(ctx, "refresh:"+userEmail).Err()
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
	// Get user_id from JWT token
	userID := c.Get("user_id").(uint)
	log.Println("User ID from token:", userID)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	// Get All data from request body
	userPayload := new(dto.UpdateUserRequest)
	err := c.Bind(userPayload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload : " + err.Error(),
		})
	}

	// Capitalize the first letter of the Name
	nameStr := strings.ToTitle(userPayload.Name)

	// Set default value for PhoneNumber if not provided
	if userPayload.PhoneNumber == "" {
		userPayload.PhoneNumber = "XXXXXX"
	}

	// Call the service to update the user
	_, err = u.UserService.UpdateUser(userID, models.User{
		Name:        nameStr,
		Password:    userPayload.Password,
		Role:        userPayload.Role,
		PhoneNumber: userPayload.PhoneNumber,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to update user : " + err.Error(),
		})
	}

	// Invalidate the user's tokens by logging them out
	u.Logout(c)

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User updated successfully",
	}

	return c.JSON(http.StatusOK, apiResponse)
}

func (u *UserControllerImpl) DeleteUser(c echo.Context) error {
	// Get user_id from JWT token
	userID := c.Get("user_id").(uint)
	log.Println("User ID from token:", userID)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	// Call the service to delete the user
	_, err := u.UserService.DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to delete user : " + err.Error(),
		})
	}

	// Invalidate the user's tokens by logging them out
	u.Logout(c)

	// Return success response
	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
	}

	return c.JSON(http.StatusOK, apiResponse)
}

func (u *UserControllerImpl) GetUserByID(c echo.Context) error {
	// Get user_id from JWT token
	userID := c.Get("user_id").(uint)
	log.Println("User ID from token:", userID)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "User ID is required",
		})
	}

	result, err := u.UserService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ApiResponse{
			Status:  http.StatusNotFound,
			Message: "User not found : " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "User retrieved successfully",
		Data:    result,
	}

	return c.JSON(http.StatusOK, apiResponse)
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

	// Store new access token in Redis
	

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
