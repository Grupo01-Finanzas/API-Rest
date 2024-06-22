package repository

import (
	"ApiRestFinance/internal/model/entities"
	"gorm.io/gorm"
)

// ClientRepository defines operations for managing Client entities.
type ClientRepository interface {
	CreateClient(client *entities.Client) error
	GetClientByID(clientID uint) (*entities.Client, error)
	UpdateClient(client *entities.Client) error
	DeleteClient(clientID uint) error
	// Add other necessary methods (like GetClientsByEstablishmentID if needed)
	CreateClientInTransaction(tx *gorm.DB, client *entities.Client) error
	DeleteClientInTransaction(tx *gorm.DB, clientID uint) error
}

type clientRepository struct {
	db *gorm.DB
}

// NewClientRepository creates a new instance of the client repository.
func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

// CreateClient creates a new client in the database.
func (r *clientRepository) CreateClient(client *entities.Client) error {
	return r.db.Create(client).Error
}

// GetClientByID retrieves a client by their ID.
func (r *clientRepository) GetClientByID(clientID uint) (*entities.Client, error) {
	var client entities.Client
	err := r.db.Preload("User").First(&client, clientID).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// UpdateClient updates an existing client in the database.
func (r *clientRepository) UpdateClient(client *entities.Client) error {
	return r.db.Save(client).Error
}

// DeleteClient deletes a client from the database.
func (r *clientRepository) DeleteClient(clientID uint) error {
	return r.db.Delete(&entities.Client{}, clientID).Error
}

// CreateClientInTransaction creates a new client within a database transaction.
func (r *clientRepository) CreateClientInTransaction(tx *gorm.DB, client *entities.Client) error {
    return tx.Create(client).Error
}

// DeleteClientInTransaction deletes a client within a database transaction.
func (r *clientRepository) DeleteClientInTransaction(tx *gorm.DB, clientID uint) error {
    return tx.Delete(&entities.Client{}, clientID).Error
}
