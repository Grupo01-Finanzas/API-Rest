package request

type CreateUserRequest struct {
	DNI      string `json:"dni" binding:"required,min=8,max=8"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required,min=5"`
	Phone    string `json:"phone" binding:"required,min=9,max=9"`
}
