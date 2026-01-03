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

type CategoryController interface {
	CreateCategory(c echo.Context) error
	GetCategoryByID(c echo.Context) error
	UpdateCategory(c echo.Context) error
	DeleteCategory(c echo.Context) error
	GetAllCategories(c echo.Context) error
}

type categoryControllerImpl struct {
	CategoryService services.CategoryService
}

// GetAllCategories implements CategoryController.
func (r *categoryControllerImpl) GetAllCategories(c echo.Context) error {
	categories, err := r.CategoryService.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get categories: " + err.Error(),
		})
	}

	apiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Categories retrieved successfully",
		Data:    categories,
	}

	return c.JSON(http.StatusOK, apiResponse)
}

// CreateCategory implements CategoryController.
func (r *categoryControllerImpl) CreateCategory(c echo.Context) error {
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

	userPayload := new(dto.CategoryRequest)
	err := c.Bind(userPayload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload : " + err.Error(),
		})
	}

	categoriesName := strings.ToTitle(userPayload.CategoriesName)
	result, err := r.CategoryService.CreateCategory(models.Category{
		CategoriesName: categoriesName,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create category: " + err.Error(),
		})
	}

	ApiResponse := dto.ApiResponse{
		Status:  http.StatusCreated,
		Message: "Category created successfully",
		Data:    result,
	}

	return c.JSON(http.StatusCreated, ApiResponse)

}

// DeleteCategory implements CategoryController.
func (r *categoryControllerImpl) DeleteCategory(c echo.Context) error {
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
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid category ID: " + err.Error(),
		})
	}

	err = r.CategoryService.DeleteCategory(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete category: " + err.Error(),
		})
	}

	ApiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Category deleted successfully",
	}

	return c.JSON(http.StatusOK, ApiResponse)

}

// GetCategoryByID implements CategoryController.
func (r *categoryControllerImpl) GetCategoryByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid category ID: " + err.Error(),
		})
	}

	category, err := r.CategoryService.GetCategoryByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get category: " + err.Error(),
		})
	}

	ApiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Category retrieved successfully",
		Data:    category,
	}

	return c.JSON(http.StatusOK, ApiResponse)
}

// UpdateCategory implements CategoryController.
func (r *categoryControllerImpl) UpdateCategory(c echo.Context) error {
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

	userPayload := new(dto.CategoryRequest)
	err := c.Bind(userPayload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload : " + err.Error(),
		})
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid category ID: " + err.Error(),
		})
	}

	categoriesName := strings.ToTitle(userPayload.CategoriesName)
	result, err := r.CategoryService.UpdateCategory(uint(id), models.Category{
		CategoriesName: categoriesName,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update category: " + err.Error(),
		})
	}

	ApiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Category updated successfully",
		Data:    result,
	}

	return c.JSON(http.StatusOK, ApiResponse)
}

func NewCategoryController(db *gorm.DB) CategoryController {
	service := services.NewCategoryService(repositories.NewCategoryRepository(db))
	return &categoryControllerImpl{
		CategoryService: service,
	}
}
