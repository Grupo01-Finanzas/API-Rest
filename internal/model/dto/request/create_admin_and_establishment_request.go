package request

type CreateAdminAndEstablishmentRequest struct {
	// User fields
	DNI      string `json:"dni" binding:"required,min=8,max=8"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required,min=5"`
	Phone    string `json:"phone" binding:"required,min=9,max=9"`
	// Establishment fields
	EstablishmentRUC     string  `json:"establishment_ruc" binding:"required"`
	EstablishmentName    string  `json:"establishment_name" binding:"required"`
	EstablishmentPhone   string  `json:"establishment_phone" binding:"required"`
	EstablishmentAddress string  `json:"establishment_address" binding:"required"`
	LateFeePercentage   float64 `json:"late_fee_percentage" binding:"omitempty,gte=0"` // Optional, can be set later
}