package dto

type MenuRequest struct {
	MenuName    string  `json:"menu_name"`
	Price       float64 `json:"price"`
	CategoryID  uint    `json:"category_id"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
}
