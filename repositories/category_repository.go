package repositories

import (
	"coffee_shop/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	// Define category-related data access methods here
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
	FindByName(name string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
}

type categoryRepositoryImpl struct {
	DB *gorm.DB
}

// GetAllCategories implements CategoryRepository.
func (c *categoryRepositoryImpl) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	result := c.DB.Where("deleted_at is NULL").Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// FindByName implements CategoryRepository.
func (c *categoryRepositoryImpl) FindByName(name string) (*models.Category, error) {
	var category models.Category
	result := c.DB.Where("categories_name = ?", name).Where("deleted_at is NULL").First(&category)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) CreateCategory(category *models.Category) error {
	return c.DB.Create(&category).Error
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) DeleteCategory(id uint) error {
	result := c.DB.Where("id = ?", id).Where("deleted_at IS NULL").Delete(&models.Category{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// GetCategoryByID implements CategoryRepository.
func (c *categoryRepositoryImpl) GetCategoryByID(id uint) (models.Category, error) {
	var category models.Category
	result := c.DB.Where("id = ?", id).Where("deleted_at IS NULL").First(&category)
	if result.Error != nil {
		return models.Category{}, result.Error
	}
	return category, nil
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepositoryImpl) UpdateCategory(category *models.Category) error {
	return c.DB.Save(category).Error
}

// UpdateCategory implements CategoryRepository.
// func (c *categoryRepositoryImpl) UpdateCategory(id uint, category *models.Category) error {
// 	result := c.DB.Where("id = ?", id).Where("deleted_at IS NULL").Updates(category)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		DB: db,
	}
}
