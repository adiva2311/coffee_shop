package dto

import "coffee_shop/models"

type CategoryResponse struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"categories_name"`
}

func ToCategoryResponse(category *models.Category) CategoryResponse {
	return CategoryResponse{
		ID:           category.ID,
		CategoryName: category.CategoriesName,
	}
}
