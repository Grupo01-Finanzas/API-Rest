package request

type CreateClientRequest struct {
	UserID   uint `json:"user_id" binding:"required"`
	IsActive bool `json:"is_active"`
}
