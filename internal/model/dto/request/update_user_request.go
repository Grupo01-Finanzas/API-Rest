package request

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty"`      // Optional
	Address  string `json:"address" binding:"omitempty"`   // Optional
	Phone    string `json:"phone" binding:"omitempty"`     // Optional
	PhotoUrl string `json:"photo_url" binding:"omitempty"` // Optional
}
