package models

import "gorm.io/gorm"

type Menu struct {
	// ID        uint `gorm:"primarykey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	gorm.Model
	MenuName    string  `gorm:"not null" json:"menu_name"`
	Price       float64 `gorm:"not null" json:"price"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	CategoryID  uint    `gorm:"not null" json:"category_id"`
	Category    Category
}

func (Menu) TableName() string {
	return "menu"
}

