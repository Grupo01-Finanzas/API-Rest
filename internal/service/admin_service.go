package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/entities"
	"errors"
	"fmt"

	"ApiRestFinance/internal/repository"
)

type AdminService interface {
	CreateAdmin(admin *entities.Admin) error
	GetAllAdmins() ([]entities.Admin, error)
	GetAdminByID(adminID uint) (*entities.Admin, error)
	UpdateAdmin(admin *entities.Admin) error
	DeleteAdmin(adminID uint) error
	RegisterEstablishment(establishment *entities.Establishment, adminID uint) error
	GetEstablishmentByID(establishmentID uint) (*entities.Establishment, error)
	GetAdminByUserID(userID uint) (*entities.Admin, error)
}

type adminService struct {
	adminRepo         repository.AdminRepository
	establishmentRepo repository.EstablishmentRepository
	userRepo          repository.UserRepository
}

func NewAdminService(adminRepo repository.AdminRepository, establishmentRepo repository.EstablishmentRepository, userRepo repository.UserRepository) AdminService {
	return &adminService{adminRepo: adminRepo, establishmentRepo: establishmentRepo, userRepo: userRepo}
}

func (s *adminService) CreateAdmin(admin *entities.Admin) error {
	existingAdmin, _ := s.adminRepo.GetAdminByUserID(admin.UserID)
	if existingAdmin != nil {
		return errors.New("exist already an admin with this user id")
	}
	return s.adminRepo.CreateAdmin(admin)
}

func (s *adminService) GetAllAdmins() ([]entities.Admin, error) {
	return s.adminRepo.GetAllAdmins()
}

func (s *adminService) GetAdminByID(adminID uint) (*entities.Admin, error) {
	return s.adminRepo.GetAdminByID(adminID)
}

func (s *adminService) UpdateAdmin(admin *entities.Admin) error {
	return s.adminRepo.UpdateAdmin(admin)
}

func (s *adminService) DeleteAdmin(adminID uint) error {
	return s.adminRepo.DeleteAdmin(adminID)
}

func (s *adminService) RegisterEstablishment(establishment *entities.Establishment, adminID uint) error {
	existingAdmin, err := s.adminRepo.GetAdminByUserID(adminID)
	if err != nil {
		return fmt.Errorf("error al buscar el administrador: %w", err)
	}

	if existingAdmin != nil {
		return errors.New("admin already has an establishment")
	}

	admin, err := s.userRepo.GetUserByID(adminID)
	if err != nil {
		return fmt.Errorf("error al buscar el administrador: %w", err)
	}

	if admin == nil {
		return errors.New("admin not found")
	}

	establishment.Admin = &entities.Admin{
		UserID: admin.ID,
		User:   *admin,
	}

	// Create the establishment request
	createEstablishmentRequest := request.CreateEstablishmentRequest{
		RUC:      establishment.RUC,
		Name:     establishment.Name,
		Phone:    establishment.Phone,
		Address:  establishment.Address,
		IsActive: establishment.IsActive,
		AdminID:  admin.ID,
	}

	// Use the request to create the establishment
	_, err = s.establishmentRepo.Create(createEstablishmentRequest)
	if err != nil {
		return fmt.Errorf("error al crear el establecimiento: %w", err)
	}

	return nil
}
func (s *adminService) GetEstablishmentByID(establishmentID uint) (*entities.Establishment, error) {
	// Call the correct method in the repository
	establishmentResponse, err := s.establishmentRepo.GetByID(establishmentID)
	if err != nil {
		return nil, err
	}

	// Convert the response to the desired Establishment type
	establishment := entities.Establishment{
		RUC:      establishmentResponse.RUC,
		Name:     establishmentResponse.Name,
		Phone:    establishmentResponse.Phone,
		Address:  establishmentResponse.Address,
		IsActive: establishmentResponse.IsActive,
	}

	return &establishment, nil
}
func (s *adminService) GetAdminByUserID(userID uint) (*entities.Admin, error) {
	return s.adminRepo.GetAdminByUserID(userID)
}
