package repository

import (
	"ApiRestFinance/internal/model/entities"
	"gorm.io/gorm"
)

type ClientRepository interface {
	CreateClient(client *entities.Client) error
	GetClientByUserID(userID uint) (*entities.Client, error)
	GetAllClients() ([]entities.Client, error)
	GetClientByID(clientID uint) (*entities.Client, error)
	UpdateClient(client *entities.Client) error
	DeleteClient(clientID uint) error
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) CreateClient(client *entities.Client) error {
	return r.db.Create(client).Error
}

func (r *clientRepository) GetClientByUserID(userID uint) (*entities.Client, error) {
	var client entities.Client
	err := r.db.Where("user_id = ?", userID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) GetAllClients() ([]entities.Client, error) {
	var clients []entities.Client
	err := r.db.Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *clientRepository) GetClientByID(clientID uint) (*entities.Client, error) {
	var client entities.Client
	err := r.db.First(&client, clientID).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) UpdateClient(client *entities.Client) error {
	return r.db.Save(client).Error
}

func (r *clientRepository) DeleteClient(clientID uint) error {
	return r.db.Delete(&entities.Client{}, clientID).Error
}
