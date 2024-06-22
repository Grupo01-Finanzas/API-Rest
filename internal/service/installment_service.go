package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"fmt"
)

// InstallmentService handles installment-related operations.
type InstallmentService interface {
	CreateInstallment(req request.CreateInstallmentRequest) (*response.InstallmentResponse, error)
	GetInstallmentByID(id uint) (*response.InstallmentResponse, error)
	UpdateInstallment(id uint, req request.UpdateInstallmentRequest) (*response.InstallmentResponse, error)
	DeleteInstallment(id uint) error
	GetInstallmentsByCreditAccountID(creditAccountID uint) ([]response.InstallmentResponse, error)
	GetOverdueInstallments(creditAccountID uint) ([]response.InstallmentResponse, error)
}

type installmentService struct {
	installmentRepo repository.InstallmentRepository
}

// NewInstallmentService creates a new instance of InstallmentService.
func NewInstallmentService(installmentRepo repository.InstallmentRepository) InstallmentService {
	return &installmentService{installmentRepo: installmentRepo}
}

// CreateInstallment creates a new installment.
func (s *installmentService) CreateInstallment(req request.CreateInstallmentRequest) (*response.InstallmentResponse, error) {
	installment := entities.Installment{
		CreditAccountID: req.CreditAccountID,
		DueDate:         req.DueDate,
		Amount:          req.Amount,
		Status:          enums.Pending, // Assuming new installments are initially pending
	}
	// Call the correct method and pass the installment as a slice
	err := s.installmentRepo.CreateInstallments([]entities.Installment{installment}) 
	if err != nil {
		return nil, fmt.Errorf("error creating installment: %w", err)
	}
	return installmentToResponse(&installment), nil
}

// GetInstallmentByID retrieves an installment by ID.
func (s *installmentService) GetInstallmentByID(id uint) (*response.InstallmentResponse, error) {
	installment, err := s.installmentRepo.GetInstallmentByID(id)
	if err != nil {
		return nil, err
	}
	return installmentToResponse(installment), nil
}

// UpdateInstallment updates an existing installment.
func (s *installmentService) UpdateInstallment(id uint, req request.UpdateInstallmentRequest) (*response.InstallmentResponse, error) {
	installment, err := s.installmentRepo.GetInstallmentByID(id)
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

	err = s.installmentRepo.UpdateInstallment(installment)
	if err != nil {
		return nil, err
	}
	return installmentToResponse(installment), nil
}

// DeleteInstallment deletes an installment.
func (s *installmentService) DeleteInstallment(id uint) error {
	return s.installmentRepo.DeleteInstallment(id)
}

// GetInstallmentsByCreditAccountID retrieves all installments for a specific credit account.
func (s *installmentService) GetInstallmentsByCreditAccountID(creditAccountID uint) ([]response.InstallmentResponse, error) {
	installments, err := s.installmentRepo.GetInstallmentsByCreditAccountID(creditAccountID)
	if err != nil {
		return nil, err
	}

	var installmentResponses []response.InstallmentResponse
	for _, installment := range installments {
		installmentResponses = append(installmentResponses, *installmentToResponse(&installment))
	}

	return installmentResponses, nil
}

// GetOverdueInstallments retrieves all overdue installments for a specific credit account.
func (s *installmentService) GetOverdueInstallments(creditAccountID uint) ([]response.InstallmentResponse, error) {
	installments, err := s.installmentRepo.GetOverdueInstallments(creditAccountID)
	if err != nil {
		return nil, err
	}

	var installmentResponses []response.InstallmentResponse
	for _, installment := range installments {
		installmentResponses = append(installmentResponses, *installmentToResponse(&installment))
	}

	return installmentResponses, nil
}

func installmentToResponse(installment *entities.Installment) *response.InstallmentResponse {
	return &response.InstallmentResponse{
		ID:              installment.ID,
		CreditAccountID: installment.CreditAccountID,
		DueDate:         installment.DueDate,
		Amount:          installment.Amount,
		Status:          installment.Status,
		CreatedAt:       installment.CreatedAt,
		UpdatedAt:       installment.UpdatedAt,
	}
}