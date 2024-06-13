package response

import "time"

type AdminResponse struct {
	ID              uint                  `json:"id"`
	UserID          uint                  `json:"user_id"`
	EstablishmentID uint                  `json:"establishment_id"`
	IsActive        bool                  `json:"is_active"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	Establishment   EstablishmentResponse `json:"establishment"`
}
