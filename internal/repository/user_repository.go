package repository

import (
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(userID uint) (*entities.User, error)
	UpdateUser(user *entities.User) error
	UpdateRol(userID uint, nuevoRol enums.Role) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(userID uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateRol(userID uint, nuevoRol enums.Role) error {
	return r.db.Model(&entities.User{}).Where("id = ?", userID).Update("Rol", nuevoRol).Error
}
