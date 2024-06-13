package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/repository"
)

type LateFeeRuleService interface {
	CreateLateFeeRule(req request.CreateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error)
	GetLateFeeRuleByID(id uint) (*response.LateFeeRuleResponse, error)
	UpdateLateFeeRule(id uint, req request.UpdateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error)
	DeleteLateFeeRule(id uint) error
	GetAllLateFeeRules() ([]response.LateFeeRuleResponse, error)
	GetLateFeeRulesByEstablishmentID(establishmentID uint) ([]response.LateFeeRuleResponse, error)
}

type lateFeeRuleService struct {
	lateFeeRuleRepo repository.LateFeeRuleRepository
}

// NewLateFeeRuleService creates a new instance of LateFeeRuleService.
func NewLateFeeRuleService(lateFeeRuleRepo repository.LateFeeRuleRepository) LateFeeRuleService {
	return &lateFeeRuleService{lateFeeRuleRepo: lateFeeRuleRepo}
}

// CreateLateFeeRule creates a new late fee rule.
func (s *lateFeeRuleService) CreateLateFeeRule(req request.CreateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error) {
	return s.lateFeeRuleRepo.Create(req)
}

// GetLateFeeRuleByID retrieves a late fee rule by ID.
func (s *lateFeeRuleService) GetLateFeeRuleByID(id uint) (*response.LateFeeRuleResponse, error) {
	return s.lateFeeRuleRepo.GetByID(id)
}

// UpdateLateFeeRule updates an existing late fee rule.
func (s *lateFeeRuleService) UpdateLateFeeRule(id uint, req request.UpdateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error) {
	return s.lateFeeRuleRepo.Update(id, req)
}

// DeleteLateFeeRule deletes a late fee rule.
func (s *lateFeeRuleService) DeleteLateFeeRule(id uint) error {
	return s.lateFeeRuleRepo.Delete(id)
}

// GetAllLateFeeRules retrieves all late fee rules.
func (s *lateFeeRuleService) GetAllLateFeeRules() ([]response.LateFeeRuleResponse, error) {
	return s.lateFeeRuleRepo.GetAll()
}

// GetLateFeeRulesByEstablishmentID retrieves all late fee rules for a specific establishment.
func (s *lateFeeRuleService) GetLateFeeRulesByEstablishmentID(establishmentID uint) ([]response.LateFeeRuleResponse, error) {
	return s.lateFeeRuleRepo.GetByEstablishmentID(establishmentID)
}
