package request

type CreateAdminRequest struct {
	UserID          uint `json:"user_id" binding:"required"`
	EstablishmentID uint `json:"establishment_id" binding:"required"`
	IsActive        bool `json:"is_active"`
}
