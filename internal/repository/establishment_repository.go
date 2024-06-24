package repository

import (
	"ApiRestFinance/internal/model/entities"
	"fmt"

	"gorm.io/gorm"
)

// EstablishmentRepository defines operations for managing Establishment entities.
type EstablishmentRepository interface {
	CreateEstablishment(establishment *entities.Establishment) error
	GetEstablishmentByID(establishmentID uint) (*entities.Establishment, error)
	UpdateEstablishment(establishment *entities.Establishment) error
	DeleteEstablishment(establishmentID uint) error
	GetEstablishmentByAdminID(adminID uint) (*entities.Establishment, error)
	CreateEstablishmentInTransaction(tx *gorm.DB, establishment *entities.Establishment) error
	CreateAdminAndEstablishment(user *entities.User, establishment *entities.Establishment) error
	GetAdminByUserID(userID uint) (*entities.User, error)
}

type establishmentRepository struct {
	db *gorm.DB
}

// NewEstablishmentRepository creates a new EstablishmentRepository instance.
func NewEstablishmentRepository(db *gorm.DB) EstablishmentRepository {
	return &establishmentRepository{db: db}
}

// CreateEstablishment creates a new establishment in the database.
func (r *establishmentRepository) CreateEstablishment(establishment *entities.Establishment) error {
	return r.db.Create(establishment).Error
}

// GetEstablishmentByID retrieves an establishment by its ID.
func (r *establishmentRepository) GetEstablishmentByID(establishmentID uint) (*entities.Establishment, error) {
	var establishment entities.Establishment
	err := r.db.First(&establishment, establishmentID).Error
	if err != nil {
		return nil, err
	}
	return &establishment, nil
}

func (r *establishmentRepository) GetEstablishmentByUserID(userID uint) (*entities.Establishment, error) {
	var establishment entities.Establishment
	err := r.db.Where("admin_id = ?", userID).First(&establishment).Error
	if err != nil {
		return nil, err
	}
	return &establishment, nil
}

// UpdateEstablishment updates an existing establishment in the database.
func (r *establishmentRepository) UpdateEstablishment(establishment *entities.Establishment) error {
	return r.db.Save(establishment).Error
}

// DeleteEstablishment deletes an establishment from the database.
func (r *establishmentRepository) DeleteEstablishment(establishmentID uint) error {
	return r.db.Delete(&entities.Establishment{}, establishmentID).Error
}

// GetEstablishmentByAdminID retrieves the establishment associated with a specific admin.
func (r *establishmentRepository) GetEstablishmentByAdminID(adminID uint) (*entities.Establishment, error) {
	var establishment entities.Establishment
	err := r.db.Where("admin_id = ?", adminID).First(&establishment).Error
	if err != nil {
		return nil, err
	}
	return &establishment, nil
}

func (r *establishmentRepository) CreateEstablishmentInTransaction(tx *gorm.DB, establishment *entities.Establishment) error {
	return tx.Create(establishment).Error
}

func (r *establishmentRepository) CreateAdminAndEstablishment(user *entities.User, establishment *entities.Establishment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		establishment.AdminID = user.ID
		if err := tx.Create(establishment).Error; err != nil {
			return fmt.Errorf("error creating establishment: %w", err)
		}

		return nil
	})
}

func (r *establishmentRepository) GetAdminByUserID(userID uint) (*entities.User, error) {
	var admin entities.User
	err := r.db.Where("id = ?", userID).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
