package dto

type CategoryRequest struct {
	CategoriesName string `json:"categories_name" binding:"required"`
}
