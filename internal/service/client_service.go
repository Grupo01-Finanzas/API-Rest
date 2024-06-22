package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// ClientService handles client-related operations.
type ClientService interface {
	CreateClient(req request.CreateClientRequest) (*response.ClientResponse, error)
	GetClientByID(clientID uint) (*response.ClientResponse, error)
	GetClientsByEstablishmentID(establishmentID uint) ([]response.ClientResponse, error)
	UpdateClient(clientID uint, req request.UpdateClientRequest) (*response.ClientResponse, error)
	DeleteClient(clientID uint) error
}

type clientService struct {
	userRepo          repository.UserRepository
	creditAccountRepo repository.CreditAccountRepository
}

// NewClientService creates a new ClientService instance.
func NewClientService(userRepo repository.UserRepository, creditAccountRepo repository.CreditAccountRepository) ClientService {
	return &clientService{userRepo: userRepo, creditAccountRepo: creditAccountRepo}
}

// CreateClient creates a new client user and their associated credit account.
func (s *clientService) CreateClient(req request.CreateClientRequest) (*response.ClientResponse, error) {
	user := &entities.User{
		DNI:       req.DNI,
		Email:     req.Email,
		Password:  string(generateRandomPassword()),
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Rol:       enums.CLIENT, 
	}

	creditAccount := &entities.CreditAccount{
		EstablishmentID:         req.EstablishmentID,
		ClientID:                user.ID, // Use User ID directly
		CreditLimit:             req.CreditLimit,
		MonthlyDueDate:          req.MonthlyDueDate,
		InterestRate:            req.InterestRate,
		InterestType:            req.InterestType,
		CreditType:              req.CreditType,
		GracePeriod:             req.GracePeriod,
		IsBlocked:               false,
		LastInterestAccrualDate: time.Now(),
		CurrentBalance:          0.0,
		LateFeePercentage:       req.LateFeePercentage,
	}

	// Let CreditAccountRepository handle the transaction
	if err := s.creditAccountRepo.CreateClientAndCreditAccount(user, creditAccount); err != nil {
        return nil, fmt.Errorf("error during client creation: %w", err)
    }

	return &response.ClientResponse{
		ID:        user.ID, 
		User:      NewUserResponse(user),
		IsActive:  true, 
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetClientByID retrieves a client by ID.
func (s *clientService) GetClientByID(userID uint) (*response.ClientResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client: %w", err)
	}
	return userToClientResponse(user), nil
}

// GetClientsByEstablishmentID retrieves all clients associated with an establishment.
func (s *clientService) GetClientsByEstablishmentID(establishmentID uint) ([]response.ClientResponse, error) {
	users, err := s.userRepo.GetClientsByEstablishmentID(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving clients: %w", err)
	}

	var clientResponses []response.ClientResponse
	for _, user := range users {
		clientResponses = append(clientResponses, *userToClientResponse(&user))
	}

	return clientResponses, nil
}

// UpdateClient updates an existing client.
func (s *clientService) UpdateClient(userID uint, req request.UpdateClientRequest) (*response.ClientResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client: %w", err)
	}

	// Update User fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	// Update credit account details (if provided in the request)
	if req.CreditLimit != 0 || req.MonthlyDueDate != 0 || req.InterestRate != 0 ||
		req.InterestType != "" || req.CreditType != "" || req.GracePeriod != 0 {

		creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(userID) // Use userID here
		if err != nil {
			return nil, fmt.Errorf("error retrieving credit account: %w", err)
		}

		if req.CreditLimit > 0 {
			creditAccount.CreditLimit = req.CreditLimit
		}
		if req.MonthlyDueDate > 0 {
			creditAccount.MonthlyDueDate = req.MonthlyDueDate
		}
		if req.InterestRate > 0 {
			creditAccount.InterestRate = req.InterestRate
		}
		if req.InterestType != "" {
			creditAccount.InterestType = req.InterestType
		}
		if req.CreditType != "" {
			creditAccount.CreditType = req.CreditType
		}
		if req.GracePeriod >= 0 {
			creditAccount.GracePeriod = req.GracePeriod
		}

		if err := s.creditAccountRepo.UpdateCreditAccount(creditAccount); err != nil {
			return nil, fmt.Errorf("error updating credit account: %w", err)
		}
	}
	return userToClientResponse(user), nil
}

// DeleteClient deletes a client user and their associated data.
func (s *clientService) DeleteClient(userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	// Ensure the user is a client
	if user.Rol != enums.CLIENT {
		return errors.New("user is not a client")
	}

	if err := s.creditAccountRepo.DeleteClientAndCreditAccount(userID); err != nil {
        return fmt.Errorf("error deleting client: %w", err)
    }

    return nil 
}

// generateRandomPassword generates a random password.
func generateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	password := make([]byte, 8)
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))]
	}
	return string(password)
}

// userToClientResponse converts a User entity to a ClientResponse DTO.
func userToClientResponse(user *entities.User) *response.ClientResponse {
	return &response.ClientResponse{
		ID:        user.ID,
		User:      NewUserResponse(user),
		IsActive:  true,               // You'll need to determine how you're managing client active status
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// NewUserResponse converts a User entity to a UserResponse DTO.
func NewUserResponse(user *entities.User) *response.UserResponse {
	if user == nil {
		return nil
	}
	return &response.UserResponse{
		ID:        user.ID,
		DNI:       user.DNI,
		Email:     user.Email,
		Name:      user.Name,
		Address:   user.Address,
		Phone:     user.Phone,
		PhotoUrl:  user.PhotoUrl,
		Rol:       user.Rol,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}