package repository

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"gorm.io/gorm"
	"time"
)

// LateFeeRepository defines the interface for late fee repository operations.
type LateFeeRepository interface {
	Create(req request.CreateLateFeeRequest) (*response.LateFeeResponse, error)
	GetByID(id uint) (*response.LateFeeResponse, error)
	Update(id uint, req request.UpdateLateFeeRequest) (*response.LateFeeResponse, error)
	Delete(id uint) error
	GetByCreditAccountID(creditAccountID uint) ([]response.LateFeeResponse, error)
}

type lateFeeRepository struct {
	db *gorm.DB
}

// NewLateFeeRepository creates a new instance of lateFeeRepository.
func NewLateFeeRepository(db *gorm.DB) LateFeeRepository {
	return &lateFeeRepository{db: db}
}

// Create creates a new late fee record.
func (r *lateFeeRepository) Create(req request.CreateLateFeeRequest) (*response.LateFeeResponse, error) {
	lateFee := entities.LateFee{
		CreditAccountID: req.CreditAccountID,
		Amount:          req.Amount,
		AppliedDate:     time.Now(),
	}

	err := r.db.Create(&lateFee).Error
	if err != nil {
		return nil, err
	}

	return getLateFeeResponse(&lateFee), nil
}

// GetByID retrieves a late fee by ID.
func (r *lateFeeRepository) GetByID(id uint) (*response.LateFeeResponse, error) {
	var lateFee entities.LateFee
	err := r.db.First(&lateFee, id).Error
	if err != nil {
		return nil, err
	}

	return getLateFeeResponse(&lateFee), nil
}

// Update updates an existing late fee.
func (r *lateFeeRepository) Update(id uint, req request.UpdateLateFeeRequest) (*response.LateFeeResponse, error) {
	var lateFee entities.LateFee
	err := r.db.First(&lateFee, id).Error
	if err != nil {
		return nil, err
	}

	if req.Amount > 0 {
		lateFee.Amount = req.Amount
	}

	err = r.db.Save(&lateFee).Error
	if err != nil {
		return nil, err
	}

	return getLateFeeResponse(&lateFee), nil
}

// Delete deletes a late fee.
func (r *lateFeeRepository) Delete(id uint) error {
	var lateFee entities.LateFee
	err := r.db.First(&lateFee, id).Error
	if err != nil {
		return err
	}

	return r.db.Delete(&lateFee).Error
}

// GetByCreditAccountID retrieves all late fees for a specific credit account.
func (r *lateFeeRepository) GetByCreditAccountID(creditAccountID uint) ([]response.LateFeeResponse, error) {
	var lateFees []entities.LateFee
	err := r.db.Where("credit_account_id = ?", creditAccountID).Find(&lateFees).Error
	if err != nil {
		return nil, err
	}

	var lateFeeResponses []response.LateFeeResponse
	for _, lateFee := range lateFees {
		lateFeeResponses = append(lateFeeResponses, *getLateFeeResponse(&lateFee))
	}

	return lateFeeResponses, nil
}

func getLateFeeResponse(lateFee *entities.LateFee) *response.LateFeeResponse {
	return &response.LateFeeResponse{
		ID:              lateFee.ID,
		CreditAccountID: lateFee.CreditAccountID,
		Amount:          lateFee.Amount,
		AppliedDate:     lateFee.AppliedDate,
	}
}
