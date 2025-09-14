package models

import (
	"gorm.io/gorm"
)

type Order struct {
	// ID        uint `gorm:"primarykey"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt DeletedAt `gorm:"index"`
	gorm.Model
	TotalPrice float64 `gorm:"not null" json:"total_price"`
	Note       string  `json:"note"`
	Status     string  `gorm:"not null;default:pending" json:"status"`
	UserID     uint    `gorm:"not null" json:"user_id"`
	User       User
}

func (Order) TableName() string {
	return "orders"
}
