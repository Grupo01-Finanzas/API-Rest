package repository

import (
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"time"

	"gorm.io/gorm"
)

// InstallmentRepository defines operations for managing Installment entities.
type InstallmentRepository interface {
	CreateInstallments(installments []entities.Installment) error // Batch create for efficiency
	GetInstallmentByID(installmentID uint) (*entities.Installment, error)
	GetInstallmentsByCreditAccountID(creditAccountID uint) ([]entities.Installment, error)
	UpdateInstallment(installment *entities.Installment) error
	DeleteInstallment(installmentID uint) error
	GetOverdueInstallments(creditAccountID uint) ([]entities.Installment, error)
}

type installmentRepository struct {
	db *gorm.DB
}

// NewInstallmentRepository creates a new InstallmentRepository instance.
func NewInstallmentRepository(db *gorm.DB) InstallmentRepository {
	return &installmentRepository{db: db}
}

// CreateInstallments creates multiple installments in a single database transaction.
func (r *installmentRepository) CreateInstallments(installments []entities.Installment) error {
	return r.db.Create(&installments).Error 
}

// GetInstallmentByID retrieves an installment by its ID.
func (r *installmentRepository) GetInstallmentByID(installmentID uint) (*entities.Installment, error) {
	var installment entities.Installment
	err := r.db.First(&installment, installmentID).Error
	if err != nil {
		return nil, err
	}
	return &installment, nil
}

// GetInstallmentsByCreditAccountID retrieves all installments for a specific credit account.
func (r *installmentRepository) GetInstallmentsByCreditAccountID(creditAccountID uint) ([]entities.Installment, error) {
	var installments []entities.Installment
	err := r.db.Where("credit_account_id = ?", creditAccountID).Find(&installments).Error
	if err != nil {
		return nil, err
	}
	return installments, nil
}

// UpdateInstallment updates an existing installment in the database.
func (r *installmentRepository) UpdateInstallment(installment *entities.Installment) error {
	return r.db.Save(installment).Error
}

// DeleteInstallment deletes an installment from the database.
func (r *installmentRepository) DeleteInstallment(installmentID uint) error {
	return r.db.Delete(&entities.Installment{}, installmentID).Error
}

// GetOverdueInstallments retrieves overdue installments for a credit account.
func (r *installmentRepository) GetOverdueInstallments(creditAccountID uint) ([]entities.Installment, error) {
	var overdueInstallments []entities.Installment
	err := r.db.Where("credit_account_id = ? AND due_date < ? AND status = ?", creditAccountID, time.Now(), enums.Overdue).Find(&overdueInstallments).Error
	if err != nil {
		return nil, err
	}
	return overdueInstallments, nil
}