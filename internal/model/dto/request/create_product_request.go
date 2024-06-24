package request

type CreateProductRequest struct {
	EstablishmentID uint    `json:"establishment_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Category        string  `json:"category" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	Price           float64 `json:"price" binding:"required,gt=0.0"`
	Stock           int     `json:"stock" binding:"required,gte=0"`
	ImageUrl        string  `json:"image_url" binding:"omitempty"`
	IsActive        bool    `json:"is_active"`
}
