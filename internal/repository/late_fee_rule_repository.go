package repository

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"gorm.io/gorm"
)

// LateFeeRuleRepository defines the interface for late fee rule repository operations.
type LateFeeRuleRepository interface {
	Create(req request.CreateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error)
	GetByID(id uint) (*response.LateFeeRuleResponse, error)
	Update(id uint, req request.UpdateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error)
	Delete(id uint) error
	GetAll() ([]response.LateFeeRuleResponse, error)
	GetByEstablishmentID(establishmentID uint) ([]response.LateFeeRuleResponse, error)
}

type lateFeeRuleRepository struct {
	db *gorm.DB
}

// NewLateFeeRuleRepository creates a new instance of lateFeeRuleRepository.
func NewLateFeeRuleRepository(db *gorm.DB) LateFeeRuleRepository {
	return &lateFeeRuleRepository{db: db}
}

// Create creates a new late fee rule.
func (r *lateFeeRuleRepository) Create(req request.CreateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error) {
	lateFeeRule := entities.LateFeeRule{
		EstablishmentID: req.EstablishmentID,
		Name:            req.Name,
		DaysOverdueMin:  req.DaysOverdueMin,
		DaysOverdueMax:  req.DaysOverdueMax,
		FeeType:         req.FeeType,
		FeeValue:        req.FeeValue,
	}

	err := r.db.Create(&lateFeeRule).Error
	if err != nil {
		return nil, err
	}

	return getLateFeeRuleResponse(&lateFeeRule), nil
}

// GetByID retrieves a late fee rule by ID.
func (r *lateFeeRuleRepository) GetByID(id uint) (*response.LateFeeRuleResponse, error) {
	var lateFeeRule entities.LateFeeRule
	err := r.db.First(&lateFeeRule, id).Error
	if err != nil {
		return nil, err
	}

	return getLateFeeRuleResponse(&lateFeeRule), nil
}

// Update updates an existing late fee rule.
func (r *lateFeeRuleRepository) Update(id uint, req request.UpdateLateFeeRuleRequest) (*response.LateFeeRuleResponse, error) {
	var lateFeeRule entities.LateFeeRule
	err := r.db.First(&lateFeeRule, id).Error
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		lateFeeRule.Name = req.Name
	}
	if req.DaysOverdueMin != 0 {
		lateFeeRule.DaysOverdueMin = req.DaysOverdueMin
	}
	if req.DaysOverdueMax != 0 {
		lateFeeRule.DaysOverdueMax = req.DaysOverdueMax
	}
	if req.FeeType != "" {
		lateFeeRule.FeeType = req.FeeType
	}
	if req.FeeValue != 0 {
		lateFeeRule.FeeValue = req.FeeValue
	}

	err = r.db.Save(&lateFeeRule).Error
	if err != nil {
		return nil, err
	}

	return getLateFeeRuleResponse(&lateFeeRule), nil
}

// Delete deletes a late fee rule.
func (r *lateFeeRuleRepository) Delete(id uint) error {
	var lateFeeRule entities.LateFeeRule
	err := r.db.First(&lateFeeRule, id).Error
	if err != nil {
		return err
	}

	return r.db.Delete(&lateFeeRule).Error
}

// GetAll retrieves all late fee rules.
func (r *lateFeeRuleRepository) GetAll() ([]response.LateFeeRuleResponse, error) {
	var lateFeeRules []entities.LateFeeRule
	err := r.db.Find(&lateFeeRules).Error
	if err != nil {
		return nil, err
	}

	var lateFeeRuleResponses []response.LateFeeRuleResponse
	for _, rule := range lateFeeRules {
		lateFeeRuleResponses = append(lateFeeRuleResponses, *getLateFeeRuleResponse(&rule))
	}

	return lateFeeRuleResponses, nil
}

// GetByEstablishmentID retrieves all late fee rules for a specific establishment.
func (r *lateFeeRuleRepository) GetByEstablishmentID(establishmentID uint) ([]response.LateFeeRuleResponse, error) {
	var lateFeeRules []entities.LateFeeRule
	err := r.db.Where("establishment_id = ?", establishmentID).Find(&lateFeeRules).Error
	if err != nil {
		return nil, err
	}

	var lateFeeRuleResponses []response.LateFeeRuleResponse
	for _, rule := range lateFeeRules {
		lateFeeRuleResponses = append(lateFeeRuleResponses, *getLateFeeRuleResponse(&rule))
	}

	return lateFeeRuleResponses, nil
}

func getLateFeeRuleResponse(rule *entities.LateFeeRule) *response.LateFeeRuleResponse {
	return &response.LateFeeRuleResponse{
		ID:              rule.ID,
		EstablishmentID: rule.EstablishmentID,
		Name:            rule.Name,
		DaysOverdueMin:  rule.DaysOverdueMin,
		DaysOverdueMax:  rule.DaysOverdueMax,
		FeeType:         rule.FeeType,
		FeeValue:        rule.FeeValue,
	}
}
