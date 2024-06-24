package repository

import (
	"ApiRestFinance/internal/model/entities"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ClientRepository defines operations for managing Client entities (using the User entity with role CLIENT).
type ClientRepository interface {
	CreateClient(client *entities.User) error
	GetClientByID(clientID uint) (*entities.User, error)
	UpdateClient(client *entities.User) error
	DeleteClient(clientID uint) error
	CreateClientInTransaction(tx *gorm.DB, client *entities.User) error
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
func (r *clientRepository) CreateClient(client *entities.User) error {
	return r.db.Create(client).Error
}

// GetClientByID retrieves a client by their ID.
func (r *clientRepository) GetClientByID(clientID uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, clientID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("error retrieving client: %w", err)
	}
	return &user, nil
}

// UpdateClient updates an existing client in the database.
func (r *clientRepository) UpdateClient(client *entities.User) error {
	return r.db.Save(client).Error
}

// DeleteClient deletes a client from the database.
func (r *clientRepository) DeleteClient(clientID uint) error {
	return r.db.Delete(&entities.User{}, clientID).Error
}

// CreateClientInTransaction creates a new client within a database transaction.
func (r *clientRepository) CreateClientInTransaction(tx *gorm.DB, client *entities.User) error {
	return tx.Create(client).Error
}

// DeleteClientInTransaction deletes a client within a database transaction.
func (r *clientRepository) DeleteClientInTransaction(tx *gorm.DB, clientID uint) error {
	return tx.Delete(&entities.User{}, clientID).Error
}
