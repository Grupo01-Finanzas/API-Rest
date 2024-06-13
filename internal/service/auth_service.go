package service

import (
	"errors"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/repository"
	"ApiRestFinance/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(req *request.CreateUserRequest) error
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
	Refresh(refreshToken string) (*response.AuthResponse, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ResetPassword(req *request.ResetPasswordRequest, userID uint) error
}

type authService struct {
	userRepo          repository.UserRepository
	establishmentRepo repository.EstablishmentRepository
	jwtSecret         string
}

func NewAuthService(userRepo repository.UserRepository, establishmentRepo repository.EstablishmentRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, establishmentRepo: establishmentRepo, jwtSecret: jwtSecret}
}

func (s *authService) RegisterUser(req *request.CreateUserRequest) error {
	// Verificar si el usuario ya existe
	_, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("el correo electrónico ya está en uso")
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Crear usuario con rol "user" por defecto
	user := &entities.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Roles:    []entities.Role{{Name: "user"}}, // Rol por defecto
	}

	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	// Obtener el usuario por correo electrónico
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Comparar la contraseña ingresada con la contraseña hasheada
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar el token de acceso
	accessToken, err := util.GenerateAccessToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generar el token de refresco
	refreshToken, err := util.GenerateRefreshToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	authResponse := &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

func (s *authService) Refresh(refreshToken string) (*response.AuthResponse, error) {
	token, err := util.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Obtener los claims como jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token de refresco inválido")
	}

	// Obtener el ID del usuario de los claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("token de refresco inválido")
	}

	// Verificar si el usuario existe (opcional, pero recomendado)
	_, err = s.userRepo.GetUserByID(uint(userID))
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Generar un nuevo token de acceso
	newAccessToken, err := util.GenerateAccessToken(uint(userID), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	authResponse := &response.AuthResponse{
		AccessToken: newAccessToken,
	}
	return authResponse, nil
}

func (s *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := util.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Extract claims and assert the type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid JWT token")
	}

	return claims, nil
}

func (s *authService) ResetPassword(req *request.ResetPasswordRequest, userID uint) error {
	// Obtener el usuario de la base de datos
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	// Hashear la nueva contraseña
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Actualizar la contraseña en el modelo
	user.Password = string(newPasswordHash)

	// Guardar el usuario actualizado en la base de datos
	return s.userRepo.UpdateUser(user)
}
