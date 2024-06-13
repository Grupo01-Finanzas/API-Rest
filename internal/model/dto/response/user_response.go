package response

import (
	"time"
)

type UserResponse struct {
	ID        uint            `json:"id"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Roles     []string        `json:"roles"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Client    *ClientResponse `json:"client"`
	Admin     *AdminResponse  `json:"admin"`
}
