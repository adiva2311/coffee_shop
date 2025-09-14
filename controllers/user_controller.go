package controllers

import (
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController interface {
	// Define user-related controller methods here
	Register(c echo.Context) error
	Login(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetUserByID(c echo.Context) error
}

type UserControllerImpl struct {
	UserService services.UserService
}

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
	panic("not implemented") // TODO: Implement
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

func NewUserController(db *gorm.DB) UserControllerImpl {
	service := services.NewUserService(db)
	return UserControllerImpl{
		UserService: service,
	}
}
