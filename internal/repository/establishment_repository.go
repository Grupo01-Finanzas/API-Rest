package repository

import (
	"errors"

	"ApiRestFinance/internal/model/entities"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"gorm.io/gorm"
)

type EstablishmentRepository interface {
	GetAll() ([]response.EstablishmentResponse, error)
	GetByID(id uint) (*response.EstablishmentResponse, error)
	GetByEstablishmentID(establishmentID uint) (*entities.Establishment, error)
	Create(req *request.CreateEstablishmentRequest) (*response.EstablishmentResponse, error)
	Update(id uint, req request.UpdateEstablishmentRequest) (*response.EstablishmentResponse, error)
	Delete(id uint) error
	RegisterProducts(establishmentID uint, productIDs []uint) error
	AddClientToEstablishment(establishmentID uint, clientID uint) error
	GetLateFeeRules(establishmentID uint) ([]entities.LateFeeRule, error)
}

type establishmentRepository struct {
	db *gorm.DB
}

func NewEstablishmentRepository(db *gorm.DB) EstablishmentRepository {
	return &establishmentRepository{db: db}
}

func (r *establishmentRepository) GetAll() ([]response.EstablishmentResponse, error) {
	var establishments []entities.Establishment
	err := r.db.Preload("Admin.Establishment").Preload("Products").Find(&establishments).Error
	if err != nil {
		return nil, err
	}

	var establishmentResponses []response.EstablishmentResponse
	for _, establishment := range establishments {
		establishmentResponses = append(establishmentResponses, response.EstablishmentResponse{
			ID:        establishment.ID,
			RUC:       establishment.RUC,
			Name:      establishment.Name,
			Phone:     establishment.Phone,
			Address:   establishment.Address,
			IsActive:  establishment.IsActive,
			CreatedAt: establishment.CreatedAt,
			UpdatedAt: establishment.UpdatedAt,
			Admin:     getAdminResponse(establishment.Admin),
			Products:  getEstablishmentProductsResponse(establishment.Products),
		})
	}

	return establishmentResponses, nil
}

func (r *establishmentRepository) GetByID(id uint) (*response.EstablishmentResponse, error) {
	var establishment entities.Establishment
	err := r.db.Preload("Admin.Establishment").Preload("Products").First(&establishment, id).Error
	if err != nil {
		return nil, err
	}

	return &response.EstablishmentResponse{
		ID:        establishment.ID,
		RUC:       establishment.RUC,
		Name:      establishment.Name,
		Phone:     establishment.Phone,
		Address:   establishment.Address,
		IsActive:  establishment.IsActive,
		CreatedAt: establishment.CreatedAt,
		UpdatedAt: establishment.UpdatedAt,
		Admin:     getAdminResponse(establishment.Admin),
		Products:  getEstablishmentProductsResponse(establishment.Products),
	}, nil
}

func (r *establishmentRepository) GetByEstablishmentID(establishmentID uint) (*entities.Establishment, error) {
	var establishment entities.Establishment
	err := r.db.Preload("Admin.Establishment").Preload("Products").First(&establishment, establishmentID).Error
	if err != nil {
		return nil, err
	}

	return &establishment, nil

}

func (r *establishmentRepository) Create(req *request.CreateEstablishmentRequest) (*response.EstablishmentResponse, error) {

	establishment := entities.Establishment{
		RUC:      req.RUC,
		Name:     req.Name,
		Phone:    req.Phone,
		Address:  req.Address,
		IsActive: true,
	}

	err := r.db.Create(&establishment).Error
	if err != nil {
		return nil, err
	}

	return &response.EstablishmentResponse{
		ID:        establishment.ID,
		RUC:       establishment.RUC,
		Name:      establishment.Name,
		Phone:     establishment.Phone,
		Address:   establishment.Address,
		IsActive:  establishment.IsActive,
		CreatedAt: establishment.CreatedAt,
		UpdatedAt: establishment.UpdatedAt,
		Admin:     getAdminResponse(establishment.Admin),
		Products:  []response.ProductResponse{},
	}, nil
}

func (r *establishmentRepository) Update(id uint, req request.UpdateEstablishmentRequest) (*response.EstablishmentResponse, error) {
	var establishment entities.Establishment
	if err := r.db.First(&establishment, id).Error; err != nil {
		return nil, errors.New("establishment not found")
	}

	establishment.RUC = req.RUC
	establishment.Name = req.Name
	establishment.Phone = req.Phone
	establishment.Address = req.Address
	establishment.IsActive = req.IsActive

	err := r.db.Save(&establishment).Error
	if err != nil {
		return nil, err
	}

	return &response.EstablishmentResponse{
		ID:        establishment.ID,
		RUC:       establishment.RUC,
		Name:      establishment.Name,
		Phone:     establishment.Phone,
		Address:   establishment.Address,
		IsActive:  establishment.IsActive,
		CreatedAt: establishment.CreatedAt,
		UpdatedAt: establishment.UpdatedAt,
		Admin:     getAdminResponse(establishment.Admin),
		Products:  getEstablishmentProductsResponse(establishment.Products),
	}, nil
}

func (r *establishmentRepository) Delete(id uint) error {
	var establishment entities.Establishment
	err := r.db.First(&establishment, id).Error
	if err != nil {
		return errors.New("establishment not found")
	}

	return r.db.Delete(&establishment).Error
}

func (r *establishmentRepository) RegisterProducts(establishmentID uint, productIDs []uint) error {
	var establishment entities.Establishment
	if err := r.db.First(&establishment, establishmentID).Error; err != nil {
		return errors.New("establishment not found")
	}

	var products []entities.Product
	if err := r.db.Find(&products, productIDs).Error; err != nil {
		return errors.New("products not found")
	}

	for _, product := range products {
		establishment.Products = append(establishment.Products, product)
	}

	return r.db.Save(&establishment).Error
}

func getEstablishmentProductsResponse(products []entities.Product) []response.ProductResponse {
	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, response.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return productResponses
}

func getAdminResponse(admin *entities.Admin) *response.AdminResponse {
	if admin == nil {
		return nil
	}

	return &response.AdminResponse{
		ID:              admin.ID,
		UserID:          admin.UserID,
		EstablishmentID: admin.EstablishmentID,
		IsActive:        admin.IsActive,
		CreatedAt:       admin.CreatedAt,
		UpdatedAt:       admin.UpdatedAt,
		Establishment: response.EstablishmentResponse{
			ID:        admin.Establishment.ID,
			RUC:       admin.Establishment.RUC,
			Name:      admin.Establishment.Name,
			Phone:     admin.Establishment.Phone,
			Address:   admin.Establishment.Address,
			IsActive:  admin.Establishment.IsActive,
			CreatedAt: admin.Establishment.CreatedAt,
			UpdatedAt: admin.Establishment.UpdatedAt,
		},
	}
}

func (r *establishmentRepository) AddClientToEstablishment(establishmentID uint, clientID uint) error {
	var establishment entities.Establishment
	err := r.db.First(&establishment, establishmentID).Error
	if err != nil {
		return errors.New("establishment not found")
	}

	var client entities.Client
	err = r.db.First(&client, clientID).Error
	if err != nil {
		return errors.New("client not found")
	}

	// Associate the client with the establishment
	establishment.Clients = append(establishment.Clients, client)

	return r.db.Save(&establishment).Error
}

func (r *establishmentRepository) GetLateFeeRules(establishmentID uint) ([]entities.LateFeeRule, error) {
	var lateFeeRules []entities.LateFeeRule
	err := r.db.Where("establishment_id = ?", establishmentID).Find(&lateFeeRules).Error
	if err != nil {
		return nil, err
	}
	return lateFeeRules, nil
}
