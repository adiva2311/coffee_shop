package controllers

import (
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/repositories"
	"coffee_shop/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MenuController interface {
	// Define menu-related controller methods here
	GetAllMenus(c echo.Context) error
	UpdateMenu(c echo.Context) error
	CreateMenu(c echo.Context) error
	DeleteMenu(c echo.Context) error
	GetMenuByID(c echo.Context) error
}

type MenuControllerImpl struct {
	MenuService services.MenuService
}

// CreateMenu implements MenuController.
func (m *MenuControllerImpl) CreateMenu(c echo.Context) error {
	// Check if user is admin
	userRole, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, dto.ApiResponse{
			Status:  http.StatusForbidden,
			Message: "Forbidden: Admins only",
		})
	}

	userPayload := new(dto.MenuRequest)
	if err := c.Bind(userPayload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	MenuName := strings.ToTitle(userPayload.MenuName)
	menuCreated, err := m.MenuService.CreateMenu(models.Menu{
		MenuName:    MenuName,
		Price:       userPayload.Price,
		Description: userPayload.Description,
		ImageURL:    userPayload.ImageURL,
		CategoryID:  userPayload.CategoryID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create menu: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusCreated,
		Message: "Menu created successfully",
		Data:    menuCreated,
	}
	return c.JSON(http.StatusCreated, apiResponse)
}

// DeleteMenu implements MenuController.
func (m *MenuControllerImpl) DeleteMenu(c echo.Context) error {
	// Check if user is admin
	userRole, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, dto.ApiResponse{
			Status:  http.StatusForbidden,
			Message: "Forbidden: Admins only",
		})
	}

	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid menu ID parameter: " + err.Error(),
		})
	}

	err = m.MenuService.DeleteMenu(uint(menuID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete menu: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Menu deleted successfully",
	}

	return c.JSON(http.StatusOK, apiResponse)
}

// GetAllMenus implements MenuController.
func (m *MenuControllerImpl) GetAllMenus(c echo.Context) error {
	allMenus, err := m.MenuService.GetAllMenus()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get menus: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Menus retrieved successfully",
		Data:    allMenus,
	}

	return c.JSON(http.StatusOK, apiResponse)
}

// GetMenuByID implements MenuController.
func (m *MenuControllerImpl) GetMenuByID(c echo.Context) error {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid menu ID parameter: " + err.Error(),
		})
	}

	menu, err := m.MenuService.GetMenuByID(uint(menuID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get menu by ID: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Menu retrieved successfully",
		Data:    menu,
	}

	return c.JSON(http.StatusOK, apiResponse)
}

// UpdateMenu implements MenuController.
func (m *MenuControllerImpl) UpdateMenu(c echo.Context) error {
	// Check if user is admin
	userRole, ok := c.Get("role").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ApiResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, dto.ApiResponse{
			Status:  http.StatusForbidden,
			Message: "Forbidden: Admins only",
		})
	}

	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid menu ID parameter: " + err.Error(),
		})
	}

	updatePayload := new(dto.MenuRequest)
	if err := c.Bind(updatePayload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	MenuName := strings.ToTitle(updatePayload.MenuName)
	updatedMenu, err := m.MenuService.UpdateMenu(uint(menuID), models.Menu{
		MenuName:    MenuName,
		Price:       updatePayload.Price,
		Description: updatePayload.Description,
		ImageURL:    updatePayload.ImageURL,
		CategoryID:  updatePayload.CategoryID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update menu: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Menu updated successfully",
		Data:    updatedMenu,
	}

	return c.JSON(http.StatusOK, apiResponse)
}

func NewMenuController(db *gorm.DB) MenuController {
	service := services.NewMenuService(repositories.NewMenuRepository(db))
	return &MenuControllerImpl{
		MenuService: service,
	}
}
