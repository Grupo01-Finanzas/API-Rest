package response

import (
	"time"
)

// ClientResponse represents the response data for a client.
type ClientResponse struct {
	ID        uint        `json:"id"`
	User      *UserResponse `json:"user"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}