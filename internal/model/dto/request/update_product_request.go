package request

import "ApiRestFinance/internal/model/entities/enums"

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty"`
	Category        enums.ProductCategory `json:"category" binding:"required"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
	ImageUrl    string  `json:"image_url" binding:"omitempty"`
	IsActive    bool    `json:"is_active"`
}