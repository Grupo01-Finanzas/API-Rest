package request

import "ApiRestFinance/internal/model/entities"

type CreateProductRequest struct {
	Name            string                   `json:"name" binding:"required"`
	Description     string                   `json:"description" binding:"required"`
	Price           float64                  `json:"price" binding:"required,gt=0"`
	Category        entities.ProductCategory `json:"category" binding:"required"`
	Stock           int                      `json:"stock" binding:"required,gte=0"`
	ImageUrl		string                   `json:"image_url"`
	IsActive        bool                     `json:"is_active"`
	EstablishmentID uint                     `json:"establishment_id" binding:"required"`
}
