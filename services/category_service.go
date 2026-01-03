package services

import (
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/repositories"
	"errors"
)

type CategoryService interface {
	CreateCategory(categoryRequest models.Category) (dto.CategoryResponse, error)
	GetCategoryByID(id uint) (dto.CategoryResponse, error)
	UpdateCategory(id uint, categoryReq models.Category) (dto.CategoryResponse, error)
	DeleteCategory(id uint) error
	GetAllCategories() ([]dto.CategoryResponse, error)
}

type CategoryServiceImpl struct {
	CategoryRepo repositories.CategoryRepository
}

// GetAllCategories implements CategoryService.
func (c *CategoryServiceImpl) GetAllCategories() ([]dto.CategoryResponse, error) {
	categories, err := c.CategoryRepo.GetAllCategories()
	if err != nil {
		return nil, errors.New("failed to get categories: " + err.Error())
	}

	var CategoryResponse []dto.CategoryResponse
	for _, category := range categories {
		CategoryResponse = append(CategoryResponse, dto.ToCategoryResponse(&category))
	}
	return CategoryResponse, nil
}

// CreateCategory implements CategoryService.
func (c *CategoryServiceImpl) CreateCategory(categoryRequest models.Category) (dto.CategoryResponse, error) {
	// Check if category already exists
	findCategory, err := c.CategoryRepo.FindByName(categoryRequest.CategoriesName)
	if err == nil && findCategory != nil {
		return dto.CategoryResponse{}, errors.New("category already exists: " + findCategory.CategoriesName)
	}

	categoryRequest = models.Category{
		CategoriesName: categoryRequest.CategoriesName,
	}

	err = c.CategoryRepo.CreateCategory(&categoryRequest)
	if err != nil {
		return dto.CategoryResponse{}, errors.New("failed to create category: " + err.Error())
	}
	return dto.ToCategoryResponse(&categoryRequest), nil
}

// DeleteCategory implements CategoryService.
func (c *CategoryServiceImpl) DeleteCategory(id uint) error {
	err := c.CategoryRepo.DeleteCategory(id)
	if err != nil {
		return errors.New("failed to delete category: " + err.Error())
	}
	return nil
}

// GetCategoryByID implements CategoryService.
func (c *CategoryServiceImpl) GetCategoryByID(id uint) (dto.CategoryResponse, error) {
	category, err := c.CategoryRepo.GetCategoryByID(id)
	if err != nil {
		return dto.CategoryResponse{}, errors.New("failed to get category: " + err.Error())
	}
	return dto.ToCategoryResponse(&category), nil
}

// UpdateCategory implements CategoryService.
func (c *CategoryServiceImpl) UpdateCategory(id uint, categoryReq models.Category) (dto.CategoryResponse, error) {
	// Get existing category
	category, err := c.CategoryRepo.GetCategoryByID(id)
	if err != nil {
		return dto.CategoryResponse{}, errors.New("category not found: " + err.Error())
	}

	// Check if category already exists
	findCategory, err := c.CategoryRepo.FindByName(categoryReq.CategoriesName)
	if err == nil && findCategory != nil {
		return dto.CategoryResponse{}, errors.New("category already exists: " + findCategory.CategoriesName)
	}

	// Update category fields
	category.CategoriesName = categoryReq.CategoriesName

	// Save updated category
	err = c.CategoryRepo.UpdateCategory(&category)
	if err != nil {
		return dto.CategoryResponse{}, errors.New("failed to update category: " + err.Error())
	}
	return dto.ToCategoryResponse(&category), nil
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepo: categoryRepo,
	}
}
