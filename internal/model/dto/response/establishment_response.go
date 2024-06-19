package response

import (
	"time"
)

type EstablishmentResponse struct {
	ID        uint              `json:"id"`
	RUC       string            `json:"ruc"`
	Name      string            `json:"name"`
	Phone     string            `json:"phone"`
	Address   string            `json:"address"`
	ImageUrl  string	        `json:"image_url"`
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Admin     *AdminResponse    `json:"admin"`
	Products  []ProductResponse `json:"products"`
}
