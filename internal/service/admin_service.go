package service

import (
	"ApiRestFinance/internal/model/dto/request"

	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"

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
	RegisterEstablishment(establishment *request.CreateEstablishmentRequest, adminID uint) error
	GetEstablishmentByAdminID(adminID uint) (*entities.Establishment, error)
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

func (s *adminService) RegisterEstablishment(establishment *request.CreateEstablishmentRequest, adminID uint) error {

	existingAdmin, err := s.adminRepo.GetAdminByUserID(adminID)

	if existingAdmin != nil {
		return errors.New("admin already has an establishment")
	}

	admin, err := s.userRepo.GetUserByID(adminID)
	if err != nil {
		return fmt.Errorf("error al buscar el usuario administrador: %w", err)
	}

	if admin == nil {
		return errors.New("admin not found")
	}

	establishmentCreated, err := s.establishmentRepo.Create(establishment)
	if err != nil {
		return fmt.Errorf("error al crear el establecimiento: %w", err)
	}

	admin.Rol = enums.ADMIN
	if err := s.userRepo.UpdateUser(admin); err != nil {
		return fmt.Errorf("error al actualizar el rol del administrador: %w", err)
	}

	establishmentEntity, err := s.establishmentRepo.GetByEstablishmentID(establishmentCreated.ID)

	newAdmin := &entities.Admin{
		UserID:          admin.ID,
		User:            admin,
		EstablishmentID: establishmentCreated.ID,
		Establishment:   establishmentEntity,
		IsActive:        true,
	}

	if err := s.adminRepo.CreateAdmin(newAdmin); err != nil {
		return fmt.Errorf("error al crear el registro del administrador: %w", err)
	}

	establishmentEntity.Admin = newAdmin

	return nil
}

func (s *adminService) GetEstablishmentByAdminID(adminID uint) (*entities.Establishment, error) {
	// Call the correct method in the repository
	establishmentResponse, err := s.adminRepo.GetEstablishmentByAdminID(adminID)
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
