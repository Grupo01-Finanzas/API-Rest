package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/repository"
)

// InstallmentService defines the interface for installment service operations.
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
	return s.installmentRepo.Create(req)
}

// GetInstallmentByID retrieves an installment by ID.
func (s *installmentService) GetInstallmentByID(id uint) (*response.InstallmentResponse, error) {
	return s.installmentRepo.GetByID(id)
}

// UpdateInstallment updates an existing installment.
func (s *installmentService) UpdateInstallment(id uint, req request.UpdateInstallmentRequest) (*response.InstallmentResponse, error) {
	return s.installmentRepo.Update(id, req)
}

// DeleteInstallment deletes an installment.
func (s *installmentService) DeleteInstallment(id uint) error {
	return s.installmentRepo.Delete(id)
}

// GetInstallmentsByCreditAccountID retrieves all installments for a specific credit account.
func (s *installmentService) GetInstallmentsByCreditAccountID(creditAccountID uint) ([]response.InstallmentResponse, error) {
	return s.installmentRepo.GetByCreditAccountID(creditAccountID)
}

// GetOverdueInstallments retrieves all overdue installments for a specific credit account.
func (s *installmentService) GetOverdueInstallments(creditAccountID uint) ([]response.InstallmentResponse, error) {
	return s.installmentRepo.GetOverdueInstallments(creditAccountID)
}
