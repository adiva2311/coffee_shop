package repositories

import (
	"coffee_shop/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Define user-related data access methods here
	RegisterUser(user *models.User) error
	CheckEmailValid(email string) (*models.User, error)
	CheckEmailExists(email string) (bool, error)
	GetUserByID(user_id uint) (*models.User, error)
	UpdateUser(user_id uint, user *models.User) error
	DeleteUser(user_id uint) (int, error)
}

type userRepository struct {
	DB *gorm.DB
}

// CheckEmailValid implements UserRepository.
func (u *userRepository) CheckEmailValid(email string) (*models.User, error) {
	var user models.User
	result := u.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CheckEmailExists implements UserRepository.
func (u *userRepository) CheckEmailExists(email string) (bool, error) {
	var count int64
	result := u.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(user_id uint) (int, error) {
	result := u.DB.Where("id = ?", user_id).Where("deleted_at IS NULL").Delete(&models.User{})
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return int(rowsAffected), nil
}

// GetUserByID implements UserRepository.
func (u *userRepository) GetUserByID(user_id uint) (*models.User, error) {
	var user models.User
	result := u.DB.Where("id = ?", user_id).Where("deleted_at IS NULL").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// RegisterUser implements UserRepository.
func (u *userRepository) RegisterUser(user *models.User) error {
	return u.DB.Create(&user).Error
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(user_id uint, user *models.User) error {
	result := u.DB.Where("id = ?", user_id).Where("deleted_at IS NULL").Updates(user)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}
