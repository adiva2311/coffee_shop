package models

import "gorm.io/gorm"

type Category struct {
	// ID        uint `gorm:"primarykey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	gorm.Model
	CategoryName string `gorm:"not null" json:"category_name"`
}

func (Category) TableName() string {
	return "categories"
}
