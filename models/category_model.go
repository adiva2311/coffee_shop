package models

import "gorm.io/gorm"

type Category struct {
	// ID        uint `gorm:"primarykey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	gorm.Model
	CategoriesName string `gorm:"not null" json:"categories_name"`
}

func (Category) TableName() string {
	return "categories"
}
