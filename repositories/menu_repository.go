package repositories

import (
	"coffee_shop/models"

	"gorm.io/gorm"
)

type MenuRepository interface {
	// Define menu-related data access methods here
	GetAllMenus() ([]models.Menu, error)
	CreateMenu(menu *models.Menu) error
	UpdateMenu(id uint, menu *models.Menu) error
	DeleteMenu(id uint) error
	GetMenuByID(id uint) (*models.Menu, error)
	FindByName(name string) (*models.Menu, error)
	FindByCategoryID(category_id uint) ([]models.Menu, error)
}

type MenuRepositoryImpl struct {
	DB *gorm.DB
}

// GetAllMenu implements MenuRepository.
func (m *MenuRepositoryImpl) GetAllMenus() ([]models.Menu, error) {
	var all_menus []models.Menu
	result := m.DB.Where("deleted_at is NULL").Find(&all_menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return all_menus, nil
}

// CreateMenu implements MenuRepository.
func (m *MenuRepositoryImpl) CreateMenu(menu *models.Menu) error {
	return m.DB.Create(&menu).Error
}

// UpdateMenu implements MenuRepository.
func (m *MenuRepositoryImpl) UpdateMenu(id uint, menu *models.Menu) error {
	result := m.DB.Where("id = ?", id).Where("deleted_at IS NULL").Updates(menu)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteMenu implements MenuRepository.
func (m *MenuRepositoryImpl) DeleteMenu(id uint) error {
	result := m.DB.Where("id = ?", id).Where("deleted_at IS NULL").Delete(&models.Menu{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// FindByCategoryID implements MenuRepository.
func (m *MenuRepositoryImpl) FindByCategoryID(category_id uint) ([]models.Menu, error) {
	var menus []models.Menu
	result := m.DB.Where("category_id = ?", category_id).Where("deleted_at is NULL").Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}

// FindByName implements MenuRepository.
func (m *MenuRepositoryImpl) FindByName(name string) (*models.Menu, error) {
	var menu models.Menu
	result := m.DB.Where("name = ?", name).Where("deleted_at is NULL").First(&menu)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menu, nil
}

// GetMenuByID implements MenuRepository.
func (m *MenuRepositoryImpl) GetMenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	result := m.DB.Where("id = ?", id).Where("deleted_at is NULL").First(&menu)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menu, nil
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &MenuRepositoryImpl{
		DB: db,
	}
}
