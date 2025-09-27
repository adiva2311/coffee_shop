package models

import "gorm.io/gorm"

type User struct {
	// ID        uint `gorm:"primarykey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Email       string `gorm:"unique;not null" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Role        string `gorm:"not null;default:customer" json:"role"` // enum('admin','cashier','customer')
	PhoneNumber string `gorm:"unique;not null" json:"phone_number"`
}

func (User) TableName() string {
	return "users"
}
