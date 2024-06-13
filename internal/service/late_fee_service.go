package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/repository"
)

// LateFeeService defines the interface for late fee service operations.
type LateFeeService interface {
	CreateLateFee(req request.CreateLateFeeRequest) (*response.LateFeeResponse, error)
	GetLateFeeByID(id uint) (*response.LateFeeResponse, error)
	UpdateLateFee(id uint, req request.UpdateLateFeeRequest) (*response.LateFeeResponse, error)
	DeleteLateFee(id uint) error
	GetLateFeesByCreditAccountID(creditAccountID uint) ([]response.LateFeeResponse, error)
}

type lateFeeService struct {
	lateFeeRepo repository.LateFeeRepository
}

// NewLateFeeService creates a new instance of LateFeeService.
func NewLateFeeService(lateFeeRepo repository.LateFeeRepository) LateFeeService {
	return &lateFeeService{lateFeeRepo: lateFeeRepo}
}

// CreateLateFee creates a new late fee record.
func (s *lateFeeService) CreateLateFee(req request.CreateLateFeeRequest) (*response.LateFeeResponse, error) {
	return s.lateFeeRepo.Create(req)
}

// GetLateFeeByID retrieves a late fee by ID.
func (s *lateFeeService) GetLateFeeByID(id uint) (*response.LateFeeResponse, error) {
	return s.lateFeeRepo.GetByID(id)
}

// UpdateLateFee updates an existing late fee.
func (s *lateFeeService) UpdateLateFee(id uint, req request.UpdateLateFeeRequest) (*response.LateFeeResponse, error) {
	return s.lateFeeRepo.Update(id, req)
}

// DeleteLateFee deletes a late fee.
func (s *lateFeeService) DeleteLateFee(id uint) error {
	return s.lateFeeRepo.Delete(id)
}

// GetLateFeesByCreditAccountID retrieves all late fees for a specific credit account.
func (s *lateFeeService) GetLateFeesByCreditAccountID(creditAccountID uint) ([]response.LateFeeResponse, error) {
	return s.lateFeeRepo.GetByCreditAccountID(creditAccountID)
}
