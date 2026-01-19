package dto

import "coffee_shop/models"

type MenuResponse struct {
	MenuName    string  `json:"menu_name"`
	Price       float64 `json:"price"`
	CategoryID  uint    `json:"category_id"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
}

func ToMenuResponse(menu *models.Menu) MenuResponse {
	return MenuResponse{
		MenuName:    menu.MenuName,
		Price:       menu.Price,
		CategoryID:  menu.CategoryID,
		Description: menu.Description,
		ImageURL:    menu.ImageURL,
	}
}
