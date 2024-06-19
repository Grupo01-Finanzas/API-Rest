package response

import (
	"ApiRestFinance/internal/model/entities"
	"time"
)

type ProductResponse struct {
	ID            uint                     `json:"id"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	Price         float64                  `json:"price"`
	Category      entities.ProductCategory `json:"category"`
	Stock         int                      `json:"stock"`
	ImageUrl	  string				   `json:"image_url"`
	IsActive      bool                     `json:"is_active"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	Establishment *EstablishmentResponse   `json:"establishment"`
}
