package service

import (
	"ApiRestFinance/internal/model/entities"
	"errors"

	"ApiRestFinance/internal/repository"
)

type ClientService interface {
	CreateClient(client *entities.Client) error
	GetAllClients() ([]entities.Client, error)
	GetClientByID(clientID uint) (*entities.Client, error)
	UpdateClient(client *entities.Client) error
	DeleteClient(clientID uint) error
}

type clientService struct {
	clientRepo repository.ClientRepository
	userRepo   repository.UserRepository // Agrega userRepo
}

func NewClientService(clientRepo repository.ClientRepository, userRepo repository.UserRepository) ClientService { // Actualiza la firma
	return &clientService{clientRepo: clientRepo, userRepo: userRepo}
}

func (s *clientService) CreateClient(client *entities.Client) error {
	// Verificar si ya existe un cliente con el mismo UserID
	existingClient, _ := s.clientRepo.GetClientByUserID(client.UserID)
	if existingClient != nil {
		return errors.New("ya existe un cliente asociado a este usuario")
	}
	return s.clientRepo.CreateClient(client)
}

func (s *clientService) GetAllClients() ([]entities.Client, error) {
	return s.clientRepo.GetAllClients()
}

func (s *clientService) GetClientByID(clientID uint) (*entities.Client, error) {
	return s.clientRepo.GetClientByID(clientID)
}

func (s *clientService) UpdateClient(client *entities.Client) error {
	return s.clientRepo.UpdateClient(client)
}

func (s *clientService) DeleteClient(clientID uint) error {
	return s.clientRepo.DeleteClient(clientID)
}
