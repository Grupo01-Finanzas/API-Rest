package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// ProductService handles product-related operations.
type ProductService interface {
	CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error)
	GetProductByID(id uint) (*response.ProductResponse, error)
	GetAllProductsByEstablishmentID(establishmentID uint) ([]response.ProductResponse, error)
	UpdateProduct(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error)
	DeleteProduct(id uint) error
	productToResponse(product *entities.Product) *response.ProductResponse
	NewEstablishmentResponseW(establishment *entities.Establishment) response.EstablishmentResponse
}

type productService struct {
	productRepo       repository.ProductRepository
	establishmentRepo repository.EstablishmentRepository
	userRepo          repository.UserRepository
}

// NewProductService creates a new ProductService instance.
func NewProductService(productRepo repository.ProductRepository, establishmentRepo repository.EstablishmentRepository, userRepo repository.UserRepository) ProductService {
	return &productService{
		productRepo:       productRepo,
		establishmentRepo: establishmentRepo,
		userRepo:          userRepo,
	}
}

// CreateProduct creates a new product for an establishment.
func (s *productService) CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error) {
	establishment, err := s.establishmentRepo.GetEstablishmentByID(req.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving establishment: %w", err)
	}
	if establishment == nil {
		return nil, fmt.Errorf("establishment with ID %d not found", req.EstablishmentID)
	}

	// Validate Category
	var validCategory bool
	for _, categoryValue := range []enums.ProductCategory{
		enums.ProductCategoryGrocery,
		enums.ProductCategoryFruitAndVeg,
		enums.ProductCategoryMeat,
		enums.ProductCategoryPoultry,
		enums.ProductCategorySeafood,
		enums.ProductCategoryBakery,
		enums.ProductCategoryLiquor,
		enums.ProductCategoryGeneralStore,
	} {
		if categoryValue == enums.ProductCategory(req.Category) {
			validCategory = true
			break
		}
	}

	if !validCategory {
		return nil, fmt.Errorf("invalid product category: %s", req.Category)
	}

	product := entities.Product{
		EstablishmentID: establishment.ID,
		Name:            req.Name,
		Category:        enums.ProductCategory(req.Category),
		Description:     req.Description,
		Price:           req.Price,
		Stock:           req.Stock,
		ImageUrl:        req.ImageUrl,
		IsActive:        true,
		Establishment:   NewEstablishment(establishment),
	}

	err = s.productRepo.CreateProduct(&product)
	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return s.productToResponse(&product), nil
}

// GetProductByID retrieves a product by its ID.
func (s *productService) GetProductByID(id uint) (*response.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	return s.productToResponse(product), nil
}

// GetAllProductsByEstablishmentID retrieves all products for a specific establishment.
func (s *productService) GetAllProductsByEstablishmentID(establishmentID uint) ([]response.ProductResponse, error) {
	products, err := s.productRepo.GetAllProductsByEstablishmentID(establishmentID)
	if err != nil {
		return nil, err
	}

	var productResponses []response.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, *s.productToResponse(&product))
	}

	return productResponses, nil
}

// UpdateProduct updates an existing product.
func (s *productService) UpdateProduct(id uint, req request.UpdateProductRequest) (*response.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Update the product fields from the request
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.ImageUrl != "" {
		product.ImageUrl = req.ImageUrl
	}
	product.IsActive = req.IsActive

	err = s.productRepo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	return s.productToResponse(product), nil
}

// DeleteProduct deletes a product.
func (s *productService) DeleteProduct(id uint) error {
	return s.productRepo.DeleteProduct(id)
}

// UploadProductImage uploads a product image and returns the URL.
func (s *productService) UploadProductImage(file *multipart.FileHeader, productID uint) (string, error) {
	// 1. File Type Validation
	allowedFileTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	isValidFileType := false

	for _, allowedType := range allowedFileTypes {
		if fileExt == allowedType {
			isValidFileType = true
			break
		}
	}

	if !isValidFileType {
		return "", ErrInvalidFileType
	}

	// 2. File Size Validation (Example: Limit to 2MB)
	if file.Size > 2*1024*1024 {
		return "", ErrFileSizeTooLarge
	}

	// 3. Create the "images_products" directory if it doesn't exist
	imagesDir := "images_products"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.Mkdir(imagesDir, 0755)
		if err != nil {
			return "", err
		}
	}

	// 4. Generate a new filename
	newFilename := fmt.Sprintf("%d%s", productID, fileExt)

	// 5. Create the full file path
	imagePath := filepath.Join(imagesDir, newFilename)

	// 6. Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening uploaded file: %w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Println("error closing uploaded file:", err)
		}
	}(src)

	// 7. Create the destination file
	dst, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("error creating image file: %w", err)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			fmt.Println("error closing destination file:", err)
		}
	}(dst)

	// 8. Copy the uploaded file contents to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("error copying image: %w", err)
	}

	// 9. Return the URL of the uploaded image
	return imagePath, nil
}

func (s *productService) productToResponse(product *entities.Product) *response.ProductResponse {
	establishment, err := s.establishmentRepo.GetEstablishmentByID(product.EstablishmentID)
	if err != nil {
		return nil
	}
	return &response.ProductResponse{
		ID:              product.ID,
		EstablishmentID: product.EstablishmentID,
		Establishment:   s.NewEstablishmentResponseW(establishment),
		Name:            product.Name,
		Category:        product.Category,
		Description:     product.Description,
		Price:           product.Price,
		Stock:           product.Stock,
		ImageUrl:        product.ImageUrl,
		IsActive:        product.IsActive,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	}
}

func (s *productService) NewEstablishmentResponseW(establishment *entities.Establishment) response.EstablishmentResponse {
	admin, err := s.userRepo.GetUserByID(establishment.AdminID)
	if err != nil {
		return response.EstablishmentResponse{}
	}

	adminResponse := &response.UserResponse{
		ID:        admin.ID,
		Email:     admin.Email,
		Name:      admin.Name,
		Phone:     admin.Phone,
		Rol:       admin.Rol,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}

	return response.EstablishmentResponse{
		ID:                establishment.ID,
		RUC:               establishment.RUC,
		Name:              establishment.Name,
		Phone:             establishment.Phone,
		Address:           establishment.Address,
		ImageUrl:          establishment.ImageUrl,
		LateFeePercentage: establishment.LateFeePercentage,
		IsActive:          establishment.IsActive,
		CreatedAt:         establishment.CreatedAt,
		UpdatedAt:         establishment.UpdatedAt,
		Admin:             adminResponse,
		AdminID:           adminResponse.ID,
	}
}

func NewEstablishment(establishment *entities.Establishment) entities.Establishment {
	return entities.Establishment{
		RUC:               establishment.RUC,
		Name:              establishment.Name,
		Phone:             establishment.Phone,
		Address:           establishment.Address,
		ImageUrl:          establishment.ImageUrl,
		LateFeePercentage: establishment.LateFeePercentage,
		IsActive:          establishment.IsActive,
		CreatedAt:         establishment.CreatedAt,
		UpdatedAt:         establishment.UpdatedAt,
	}
}
