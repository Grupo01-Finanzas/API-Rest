package response

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type UserResponse struct {
	ID        uint       `json:"id"`
	DNI       string     `json:"dni"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	PhotoUrl  string     `json:"photo_url"`
	Rol       enums.Role `json:"rol"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
