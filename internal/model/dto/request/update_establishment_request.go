package request

type UpdateEstablishmentRequest struct {
	RUC      string `json:"ruc" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
	IsActive bool   `json:"is_active"`
}
