package request

type UpdateEstablishmentRequest struct {
	RUC                 string  `json:"ruc" binding:"required"`
	Name                string  `json:"name" binding:"required"`
	Phone               string  `json:"phone" binding:"required"`
	Address             string  `json:"address" binding:"required"`
	ImageUrl            string  `json:"image_url" binding:"omitempty"`
	IsActive            bool    `json:"is_active"`
	LateFeePercentage float64 `json:"late_fee_percentage" binding:"omitempty,gte=0"` // Optional
}