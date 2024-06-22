package response

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type ProductResponse struct {
	ID            uint              `json:"id"`
	EstablishmentID uint              `json:"establishment_id"`
	Establishment   EstablishmentResponse `json:"establishment"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Category      enums.ProductCategory `json:"category"`
	Price         float64           `json:"price"`
	Stock         int               `json:"stock"`
	ImageUrl      string            `json:"image_url"`
	IsActive      bool              `json:"is_active"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}