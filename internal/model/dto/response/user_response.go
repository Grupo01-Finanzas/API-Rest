package response

import (
	"time"
)

type UserResponse struct {
	ID        uint            `json:"id"`
	DNI       string          `json:"dni"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Address   string          `json:"address"`
	Phone     string          `json:"phone"`
	Rol       string          `json:"rol"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Client    *ClientResponse `json:"client"`
	Admin     *AdminResponse  `json:"admin"`
}
