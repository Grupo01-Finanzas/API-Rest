package repository

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

// InstallmentRepository defines the interface for installment repository operations.
type InstallmentRepository interface {
	Create(req request.CreateInstallmentRequest) (*response.InstallmentResponse, error)
	GetByID(id uint) (*response.InstallmentResponse, error)
	Update(id uint, req request.UpdateInstallmentRequest) (*response.InstallmentResponse, error)
	Delete(id uint) error
	GetByCreditAccountID(creditAccountID uint) ([]response.InstallmentResponse, error)
	GetOverdueInstallments(creditAccountID uint) ([]response.InstallmentResponse, error)
}

type installmentRepository struct {
	db *gorm.DB
}

// NewInstallmentRepository creates a new instance of installmentRepository.
func NewInstallmentRepository(db *gorm.DB) InstallmentRepository {
	return &installmentRepository{db: db}
}

// Create creates a new installment.
func (r *installmentRepository) Create(req request.CreateInstallmentRequest) (*response.InstallmentResponse, error) {
	installment := entities.Installment{
		CreditAccountID: req.CreditAccountID,
		DueDate:         req.DueDate,
		Amount:          req.Amount,
		Status:          req.Status,
	}

	err := r.db.Create(&installment).Error
	if err != nil {
		return nil, err
	}

	return getInstallmentResponse(&installment), nil
}

// GetByID retrieves an installment by ID.
func (r *installmentRepository) GetByID(id uint) (*response.InstallmentResponse, error) {
	var installment entities.Installment
	err := r.db.First(&installment, id).Error
	if err != nil {
		return nil, err
	}

	return getInstallmentResponse(&installment), nil
}

// Update updates an existing installment.
func (r *installmentRepository) Update(id uint, req request.UpdateInstallmentRequest) (*response.InstallmentResponse, error) {
	var installment entities.Installment
	err := r.db.First(&installment, id).Error
	if err != nil {
		return nil, err
	}

	if !req.DueDate.IsZero() {
		installment.DueDate = req.DueDate
	}
	if req.Amount > 0 {
		installment.Amount = req.Amount
	}
	if req.Status != "" {
		installment.Status = req.Status
	}

	err = r.db.Save(&installment).Error
	if err != nil {
		return nil, err
	}

	return getInstallmentResponse(&installment), nil
}

// Delete deletes an installment.
func (r *installmentRepository) Delete(id uint) error {
	var installment entities.Installment
	err := r.db.First(&installment, id).Error
	if err != nil {
		return err
	}

	return r.db.Delete(&installment).Error
}

// GetByCreditAccountID retrieves all installments for a specific credit account.
func (r *installmentRepository) GetByCreditAccountID(creditAccountID uint) ([]response.InstallmentResponse, error) {
	var installments []entities.Installment
	err := r.db.Where("credit_account_id = ?", creditAccountID).Find(&installments).Error
	if err != nil {
		return nil, err
	}

	var installmentResponses []response.InstallmentResponse
	for _, installment := range installments {
		installmentResponses = append(installmentResponses, *getInstallmentResponse(&installment))
	}

	return installmentResponses, nil
}

// GetOverdueInstallments retrieves all overdue installments for a specific credit account.
func (r *installmentRepository) GetOverdueInstallments(creditAccountID uint) ([]response.InstallmentResponse, error) {
	var overdueInstallments []entities.Installment
	err := r.db.Where("credit_account_id = ? AND due_date < ? AND status = ?", creditAccountID, time.Now(), enums.Overdue).Find(&overdueInstallments).Error
	if err != nil {
		return nil, err
	}

	var overdueInstallmentResponses []response.InstallmentResponse
	for _, installment := range overdueInstallments {
		overdueInstallmentResponses = append(overdueInstallmentResponses, *getInstallmentResponse(&installment))
	}

	return overdueInstallmentResponses, nil
}

func getInstallmentResponse(installment *entities.Installment) *response.InstallmentResponse {
	return &response.InstallmentResponse{
		ID:              installment.ID,
		CreditAccountID: installment.CreditAccountID,
		DueDate:         installment.DueDate,
		Amount:          installment.Amount,
		Status:          installment.Status,
	}
}
