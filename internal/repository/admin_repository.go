package repository

import (
	"ApiRestFinance/internal/model/entities"
	"errors"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(admin *entities.Admin) error
	GetAdminByUserID(userID uint) (*entities.Admin, error)
	GetAllAdmins() ([]entities.Admin, error)
	GetAdminByID(adminID uint) (*entities.Admin, error)
	UpdateAdmin(admin *entities.Admin) error
	DeleteAdmin(adminID uint) error
	GetEstablishmentByAdminID(adminID uint) (*entities.Establishment, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateAdmin(admin *entities.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepository) GetAdminByUserID(userID uint) (*entities.Admin, error) {
	var admin entities.Admin
	err := r.db.Where("user_id = ?", userID).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) GetAllAdmins() ([]entities.Admin, error) {
	var admins []entities.Admin
	err := r.db.Find(&admins).Error
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *adminRepository) GetAdminByID(adminID uint) (*entities.Admin, error) {
	var admin entities.Admin
	err := r.db.First(&admin, adminID).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) UpdateAdmin(admin *entities.Admin) error {
	return r.db.Save(admin).Error
}

func (r *adminRepository) DeleteAdmin(adminID uint) error {
	return r.db.Delete(&entities.Admin{}, adminID).Error
}

func (r *adminRepository) GetEstablishmentByAdminID(adminID uint) (*entities.Establishment, error) {
	var admin entities.Admin
	err := r.db.Preload("Establishment").First(&admin, adminID).Error
	if err != nil {
		return nil, err
	}

	if admin.EstablishmentID == 0 {
		return nil, errors.New("establishment not found for admin")
	}

	establishmentEntity := admin.Establishment

	return establishmentEntity, nil
}
