package repository

import (
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository defines operations for managing User entities.
type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(userID uint) (*entities.User, error)
	UpdateUser(user *entities.User) error
	DeleteUser(userID uint) error
	CreateUserInTransaction(tx *gorm.DB, user *entities.User) error
	CreateClientInTransaction(tx *gorm.DB, client *entities.Client) error
	DeleteClientInTransaction(tx *gorm.DB, clientID uint) error
	GetClientByID(clientID uint) (*entities.Client, error)
	UpdateClient(client *entities.Client) error
	GetClientsByEstablishmentID(establishmentID uint) ([]entities.User, error)

}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user in the database.
func (r *userRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

// GetUserByEmail retrieves a user by their email address.
func (r *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *userRepository) GetUserByID(userID uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database.
func (r *userRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}

// DeleteUser deletes a user from the database.
func (r *userRepository) DeleteUser(userID uint) error {
	return r.db.Delete(&entities.User{}, userID).Error
}

func (r *userRepository) CreateUserInTransaction(tx *gorm.DB, user *entities.User) error {
    return tx.Create(user).Error 
}

// CreateClientInTransaction creates a new client within a database transaction.
func (r *userRepository) CreateClientInTransaction(tx *gorm.DB, client *entities.Client) error {
    return tx.Create(client).Error
}

// DeleteClientInTransaction deletes a client within a database transaction.
func (r *userRepository) DeleteClientInTransaction(tx *gorm.DB, clientID uint) error {
    return tx.Delete(&entities.Client{}, clientID).Error
}

// GetClientByID retrieves a client by their ID.
func (r *userRepository) GetClientByID(clientID uint) (*entities.Client, error) {
    var client entities.Client
    err := r.db.Where("user_id = ?", clientID).Preload("User").First(&client).Error 
    if err != nil {
        return nil, err
    }
    return &client, nil
}

// UpdateClient updates an existing client in the database.
func (r *userRepository) UpdateClient(client *entities.Client) error {
    return r.db.Save(client).Error
}

// GetClientsByEstablishmentID retrieves users with the CLIENT role associated with a given establishment ID.
func (r *userRepository) GetClientsByEstablishmentID(establishmentID uint) ([]entities.User, error) {
	var clients []entities.User
	err := r.db.Joins("JOIN credit_accounts ON credit_accounts.client_id = users.id").
		Where("credit_accounts.establishment_id = ? AND users.rol = ?", establishmentID, enums.CLIENT).
		Find(&clients).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving clients: %w", err)
	}
	return clients, nil
}