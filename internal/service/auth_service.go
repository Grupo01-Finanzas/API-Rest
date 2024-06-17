package service

import (
	"ApiRestFinance/internal/model/entities/enums"
	"errors"
	"time"

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
	AttemptRefresh(accessToken string) (*response.AuthResponse, error)
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

	_, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.User{
		DNI:       req.DNI,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Rol:       enums.USER,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credentials invalids")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credentials invalids")
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
