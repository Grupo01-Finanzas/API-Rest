package repository

import (
	"errors"

	"ApiRestFinance/internal/model/entities"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]response.ProductResponse, error)
	GetByID(id uint) (*response.ProductResponse, error)
	GetByEstablishmentID(establishmentID uint) ([]response.ProductResponse, error)
	Create(req request.CreateProductRequest) (*response.ProductResponse, error)
	Update(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error)
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]response.ProductResponse, error) {
	var products []entities.Product
	err := r.db.Preload("Establishment").Find(&products).Error
	if err != nil {
		return nil, err
	}

	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, response.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price,
			Category:      product.Category,
			Stock:         product.Stock,
			IsActive:      product.IsActive,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
			Establishment: getEstablishmentResponse(&product.Establishment),
		})
	}

	return productResponses, nil
}

func (r *productRepository) GetByID(id uint) (*response.ProductResponse, error) {
	var product entities.Product
	err := r.db.Preload("Establishment").First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &response.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		Category:      product.Category,
		Stock:         product.Stock,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
		Establishment: getEstablishmentResponse(&product.Establishment),
	}, nil
}

func (r *productRepository) GetByEstablishmentID(establishmentID uint) ([]response.ProductResponse, error) {
	var products []entities.Product
	err := r.db.Preload("Establishment").Where("establishment_id = ?", establishmentID).Find(&products).Error
	if err != nil {
		return nil, err
	}

	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, response.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price,
			Category:      product.Category,
			Stock:         product.Stock,
			IsActive:      product.IsActive,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
			Establishment: getEstablishmentResponse(&product.Establishment),
		})
	}

	return productResponses, nil
}

func (r *productRepository) Create(req request.CreateProductRequest) (*response.ProductResponse, error) {
	var establishment entities.Establishment
	err := r.db.First(&establishment, req.EstablishmentID).Error
	if err != nil {
		return nil, errors.New("establishment not found")
	}

	product := entities.Product{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Category:        req.Category,
		Stock:           req.Stock,
		IsActive:        req.IsActive,
		EstablishmentID: req.EstablishmentID,
		Establishment:   establishment,
	}

	err = r.db.Create(&product).Error
	if err != nil {
		return nil, err
	}

	return &response.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		Category:      product.Category,
		Stock:         product.Stock,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
		Establishment: getEstablishmentResponse(&product.Establishment),
	}, nil
}

func (r *productRepository) Update(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error) {
	var product entities.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, errors.New("product not found")
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Category = req.Category
	product.Stock = req.Stock
	product.IsActive = req.IsActive

	err = r.db.Save(&product).Error
	if err != nil {
		return nil, err
	}

	return &response.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		Category:      product.Category,
		Stock:         product.Stock,
		IsActive:      product.IsActive,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
		Establishment: getEstablishmentResponse(&product.Establishment),
	}, nil
}

func (r *productRepository) Delete(id uint) error {
	var product entities.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return errors.New("product not found")
	}

	return r.db.Delete(&product).Error
}

func getEstablishmentResponse(establishment *entities.Establishment) *response.EstablishmentResponse {
	if establishment == nil {
		return nil
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
	}
}
