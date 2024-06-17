package response

import (
	"ApiRestFinance/internal/model/entities"
	"time"
)

type ClientResponse struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	User      *entities.User `json:"user"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
