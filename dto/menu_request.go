package dto

type MenuRequest struct {
	MenuName    string  `json:"menu_name" form:"menu_name"`
	Price       float64 `json:"price" form:"price"`
	CategoryID  uint    `json:"category_id" form:"category_id"`
	Description string  `json:"description" form:"description"`
	ImageURL    string  `json:"image_url" form:"image_url"`
}
