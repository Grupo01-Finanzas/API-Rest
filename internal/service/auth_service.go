package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"ApiRestFinance/internal/util"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication and user-related operations.
type AuthService interface {
	RegisterAdmin(req *request.CreateAdminAndEstablishmentRequest) error
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
	AttemptRefresh(accessToken string) (*response.AuthResponse, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	ResetPassword(req *request.ResetPasswordRequest, userID uint) error
}

type authService struct {
	userRepo          repository.UserRepository
	establishmentRepo repository.EstablishmentRepository
	jwtSecret         string
}

// NewAuthService creates a new instance of authService.
func NewAuthService(userRepo repository.UserRepository, establishmentRepo repository.EstablishmentRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, establishmentRepo: establishmentRepo, jwtSecret: jwtSecret}
}

// RegisterAdmin registers a new admin user along with their establishment.
func (s *authService) RegisterAdmin(req *request.CreateAdminAndEstablishmentRequest) error {
	// Check if the email is already in use
	_, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create the User entity
	user := &entities.User{
		DNI:       req.DNI,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Rol:       enums.ADMIN,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create the Establishment entity
	establishment := &entities.Establishment{
		RUC:               req.EstablishmentRUC,
		Name:              req.EstablishmentName,
		Phone:             req.EstablishmentPhone,
		Address:           req.EstablishmentAddress,
		ImageUrl:          "", // You can handle default image URLs here
		LateFeePercentage: req.LateFeePercentage,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}


	if err := s.establishmentRepo.CreateAdminAndEstablishment(user, establishment); err != nil {
        return fmt.Errorf("error registering admin and establishment: %w", err)
    }

    return nil 
}

// Login authenticates a user with email and password.
func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := util.GenerateAccessToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

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

// AttemptRefresh attempts to refresh the access token using the refresh token.
func (s *authService) AttemptRefresh(accessToken string) (*response.AuthResponse, error) {
	token, err := util.ValidateToken(accessToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("access token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("access token invalid")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("access token invalid")
	}

	expirationTime := time.Unix(int64(exp), 0)
	if time.Since(expirationTime) > 5*time.Minute {
		return nil, errors.New("token expired, login again") 
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("access token invalid")
	}

	userID := uint(userIDFloat)
	newAccessToken, err := util.GenerateAccessToken(userID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	authResponse := &response.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newAccessToken, 
	}
	return authResponse, nil
}

// ValidateToken validates a JWT token.
func (s *authService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := util.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid JWT token")
	}

	return claims, nil
}

// ResetPassword resets the password for a user.
func (s *authService) ResetPassword(req *request.ResetPasswordRequest, userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password incorrect")
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(newPasswordHash)

	return s.userRepo.UpdateUser(user)
}
