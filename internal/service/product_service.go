package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"

	"ApiRestFinance/internal/repository"
)

type ProductService interface {
	GetAll() ([]response.ProductResponse, error)
	GetByID(id uint) (*response.ProductResponse, error)
	GetByEstablishmentID(establishmentID uint) ([]response.ProductResponse, error)
	Create(req request.CreateProductRequest) (*response.ProductResponse, error)
	Update(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error)
	Delete(id uint) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (s *productService) GetAll() ([]response.ProductResponse, error) {
	products, err := s.productRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) GetByID(id uint) (*response.ProductResponse, error) {
	product, err := s.productRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetByEstablishmentID(establishmentID uint) ([]response.ProductResponse, error) {
	products, err := s.productRepository.GetByEstablishmentID(establishmentID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) Create(req request.CreateProductRequest) (*response.ProductResponse, error) {
	product, err := s.productRepository.Create(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Update(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error) {
	product, err := s.productRepository.Update(id, req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Delete(id uint) error {
	return s.productRepository.Delete(id)
}
