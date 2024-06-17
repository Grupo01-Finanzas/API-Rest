package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"

	"ApiRestFinance/internal/repository"
)

type EstablishmentService interface {
	Create(req *request.CreateEstablishmentRequest) (*response.EstablishmentResponse, error)
	GetAll() ([]response.EstablishmentResponse, error)
	GetByID(id uint) (*response.EstablishmentResponse, error)
	Update(id uint, req request.UpdateEstablishmentRequest) (*response.EstablishmentResponse, error)
	Delete(id uint) error
	RegisterProducts(establishmentID uint, productIDs []uint) error
	AddClientToEstablishment(establishmentID uint, clientID uint) error
}

type establishmentService struct {
	establishmentRepository repository.EstablishmentRepository
}

func NewEstablishmentService(establishmentRepository repository.EstablishmentRepository) EstablishmentService {
	return &establishmentService{
		establishmentRepository: establishmentRepository,
	}
}

func (s *establishmentService) Create(req *request.CreateEstablishmentRequest) (*response.EstablishmentResponse, error) {
	establishment, err := s.establishmentRepository.Create(req)
	if err != nil {
		return nil, err
	}
	return establishment, nil
}

func (s *establishmentService) GetAll() ([]response.EstablishmentResponse, error) {
	establishments, err := s.establishmentRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return establishments, nil
}

func (s *establishmentService) GetByID(id uint) (*response.EstablishmentResponse, error) {
	establishment, err := s.establishmentRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return establishment, nil
}

func (s *establishmentService) Update(id uint, req request.UpdateEstablishmentRequest) (*response.EstablishmentResponse, error) {
	establishment, err := s.establishmentRepository.Update(id, req)
	if err != nil {
		return nil, err
	}
	return establishment, nil
}

func (s *establishmentService) Delete(id uint) error {
	return s.establishmentRepository.Delete(id)
}

func (s *establishmentService) RegisterProducts(establishmentID uint, productIDs []uint) error {
	return s.establishmentRepository.RegisterProducts(establishmentID, productIDs)
}

func (s *establishmentService) AddClientToEstablishment(establishmentID uint, clientID uint) error {
	return s.establishmentRepository.AddClientToEstablishment(establishmentID, clientID)
}
