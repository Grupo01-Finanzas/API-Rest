package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"errors"
	"fmt"
)

// AdminService handles operations related to admin users.
type AdminService interface {
	GetAdminByUserID(userID uint) (*response.AdminResponse, error)
	UpdateAdmin(userID uint, req request.UpdateUserRequest) (*response.AdminResponse, error)
	// You can add other admin-specific methods here if needed
}

type adminService struct {
	establishmentRepo repository.EstablishmentRepository
	userRepo          repository.UserRepository
}

// NewAdminService creates a new instance of adminService.
func NewAdminService(establishmentRepo repository.EstablishmentRepository, userRepo repository.UserRepository) AdminService {
	return &adminService{establishmentRepo: establishmentRepo, userRepo: userRepo}
}

// GetAdminByUserID retrieves admin details by user ID.
func (s *adminService) GetAdminByUserID(userID uint) (*response.AdminResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Make sure the user is an admin
	if user.Rol != enums.ADMIN {
		return nil, errors.New("user is not an admin")
	}

	establishment, err := s.establishmentRepo.GetEstablishmentByAdminID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving establishment: %w", err)
	}

	return &response.AdminResponse{
		User:          NewUserResponse(user),
		Establishment: establishmentToResponse(establishment, NewUserResponse(user)),
	}, nil
}

// UpdateAdmin updates an admin user's details.
func (s *adminService) UpdateAdmin(userID uint, req request.UpdateUserRequest) (*response.AdminResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Make sure the user is an admin
	if user.Rol != enums.ADMIN {
		return nil, errors.New("user is not an admin")
	}

	// Update user details
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.PhotoUrl != "" {
		user.PhotoUrl = req.PhotoUrl
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	// Retrieve the updated establishment (for the response)
	establishment, err := s.establishmentRepo.GetEstablishmentByAdminID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving establishment: %w", err)
	}

	return &response.AdminResponse{
		User:          NewUserResponse(user),
		Establishment: establishmentToResponse(establishment, NewUserResponse(user)),
	}, nil
}

// establishmentToResponse converts an Establishment entity to an EstablishmentResponse DTO.
func establishmentToResponse(establishment *entities.Establishment, admin *response.UserResponse) *response.EstablishmentResponse {
	if establishment == nil {
		return nil
	}
	return &response.EstablishmentResponse{
		ID:                establishment.ID,
		RUC:               establishment.RUC,
		Name:              establishment.Name,
		Phone:             establishment.Phone,
		Address:           establishment.Address,
		ImageUrl:          establishment.ImageUrl,
		LateFeePercentage: establishment.LateFeePercentage,
		IsActive:          establishment.IsActive,
		CreatedAt:         establishment.CreatedAt,
		UpdatedAt:         establishment.UpdatedAt,
		AdminID:           establishment.AdminID,
		Admin:             admin,
	}
}
