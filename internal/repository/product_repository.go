package repository

import (
	"ApiRestFinance/internal/model/entities"
	"gorm.io/gorm"
)

// ProductRepository defines operations for managing Product entities.
type ProductRepository interface {
	CreateProduct(product *entities.Product) error
	GetProductByID(productID uint) (*entities.Product, error)
	GetAllProductsByEstablishmentID(establishmentID uint) ([]entities.Product, error)
	UpdateProduct(product *entities.Product) error
	DeleteProduct(productID uint) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository instance.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// CreateProduct creates a new product in the database.
func (r *productRepository) CreateProduct(product *entities.Product) error {
	return r.db.Create(product).Error
}

// GetProductByID retrieves a product by its ID.
func (r *productRepository) GetProductByID(productID uint) (*entities.Product, error) {
	var product entities.Product
	err := r.db.First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProductsByEstablishmentID retrieves all products associated with a specific establishment.
func (r *productRepository) GetAllProductsByEstablishmentID(establishmentID uint) ([]entities.Product, error) {
	var products []entities.Product
	err := r.db.Where("establishment_id = ?", establishmentID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct updates an existing product in the database.
func (r *productRepository) UpdateProduct(product *entities.Product) error {
	return r.db.Save(product).Error
}

// DeleteProduct deletes a product from the database.
func (r *productRepository) DeleteProduct(productID uint) error {
	return r.db.Delete(&entities.Product{}, productID).Error
}