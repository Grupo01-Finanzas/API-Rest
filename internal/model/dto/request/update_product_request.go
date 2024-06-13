package request

import "ApiRestFinance/internal/model/entities"

type UpdateProductRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Description string                   `json:"description" binding:"required"`
	Price       float64                  `json:"price" binding:"required,gt=0"`
	Category    entities.ProductCategory `json:"category" binding:"required"`
	Stock       int                      `json:"stock" binding:"required,gte=0"`
	IsActive    bool                     `json:"is_active"`
}
