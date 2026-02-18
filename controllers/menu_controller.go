package controllers

import (
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/repositories"
	"coffee_shop/services"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func storeImage(c echo.Context, menuName string) (string, error) {
	// Limit file size to 2MB
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 2<<20)
	if err := c.Request().ParseMultipartForm(2 << 20); err != nil {
		return "", err
	}
	defer c.Request().MultipartForm.RemoveAll()

	// Get file from FORM
	file, err := c.FormFile("menu_image")
	if err != nil {
		return "", err
	}
	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("file size exceeds 2MB limit")
	}

	// Validate file type
	fileType := file.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/jpg" {
		return "", fmt.Errorf("only JPEG, JPG, and PNG images are allowed")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file
	storagePath := "./image/menu_img/"
	file.Filename = menuName + filepath.Ext(file.Filename)
	path := filepath.Join(storagePath, file.Filename)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the file content to destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", c.JSON(http.StatusBadRequest, err)
	}

	return storagePath + file.Filename, nil
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

	// Store Image
	imgURL, err := storeImage(c, userPayload.MenuName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Gagal menyimpan gambar: " + err.Error(),
			Data:    nil,
		})
	}

	MenuName := strings.ToTitle(userPayload.MenuName)
	menuCreated, err := m.MenuService.CreateMenu(models.Menu{
		MenuName:    MenuName,
		Price:       userPayload.Price,
		Description: userPayload.Description,
		ImageURL:    imgURL,
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

	// Store Image
	imgURL, err := storeImage(c, updatePayload.MenuName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Gagal menyimpan gambar: " + err.Error(),
			Data:    nil,
		})
	}

	MenuName := strings.ToTitle(updatePayload.MenuName)
	updatedMenu, err := m.MenuService.UpdateMenu(uint(menuID), models.Menu{
		MenuName:    MenuName,
		Price:       updatePayload.Price,
		Description: updatePayload.Description,
		ImageURL:    imgURL,
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
