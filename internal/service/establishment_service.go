package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EstablishmentService handles establishment-related operations.
type EstablishmentService interface {
	CreateEstablishment(req *request.CreateEstablishmentRequest, adminID uint) (*response.EstablishmentResponse, error)
	GetEstablishmentByAdminID(adminID uint) (*response.EstablishmentResponse, error)
	UpdateEstablishmentByAdminID(adminID uint, req request.UpdateEstablishmentRequest) (*response.EstablishmentResponse, error)
}

type establishmentService struct {
	establishmentRepo repository.EstablishmentRepository
	userRepo          repository.UserRepository
}

// NewEstablishmentService creates a new instance of establishmentService.
func NewEstablishmentService(establishmentRepo repository.EstablishmentRepository, userRepo repository.UserRepository) EstablishmentService {
	return &establishmentService{establishmentRepo: establishmentRepo, userRepo: userRepo}
}

// CreateEstablishment creates a new establishment for an admin user.
func (s *establishmentService) CreateEstablishment(req *request.CreateEstablishmentRequest, adminID uint) (*response.EstablishmentResponse, error) {
	// Check if the admin already has an establishment
	existingEstablishment, _ := s.establishmentRepo.GetEstablishmentByAdminID(adminID)
	if existingEstablishment != nil {
		return nil, fmt.Errorf("admin already has an establishment")
	}

	// Create the Establishment entity
	establishment := &entities.Establishment{
		RUC:               req.RUC,
		Name:              req.Name,
		Phone:             req.Phone,
		Address:           req.Address,
		ImageUrl:          req.ImageUrl,
		LateFeePercentage: req.LateFeePercentage,
		IsActive:          true,
		AdminID:           adminID,
	}

	admin, err := s.userRepo.GetUserByID(adminID)

	if err != nil {
		return nil, err
	}

	adminResponse := &response.UserResponse{
		ID:      admin.ID,
		DNI:     admin.DNI,
		Email:   admin.Email,
		Name:    admin.Name,
		Rol:     admin.Rol,
		Address: admin.Address,
		Phone:   admin.Phone,
	}

	if err := s.establishmentRepo.CreateEstablishment(establishment); err != nil {
		return nil, fmt.Errorf("error creating establishment: %w", err)
	}

	return establishmentToResponse(establishment, adminResponse), nil // Return the EstablishmentResponse here
}

// GetEstablishmentByAdminID retrieves the establishment associated with a specific admin.
func (s *establishmentService) GetEstablishmentByAdminID(adminID uint) (*response.EstablishmentResponse, error) {
	establishment, err := s.establishmentRepo.GetEstablishmentByAdminID(adminID)
	if err != nil {
		return nil, err
	}

	admin, err := s.userRepo.GetUserByID(establishment.AdminID)
	if err != nil {
		return nil, err
	}

	adminResponse := &response.UserResponse{
		ID:       admin.ID,
		DNI:      admin.DNI,
		Email:    admin.Email,
		Name:     admin.Name,
		Rol:      admin.Rol,
		Address:  admin.Address,
		Phone:    admin.Phone,
		PhotoUrl: admin.PhotoUrl,
	}

	// Convert to Response Type
	establishmentResponse := &response.EstablishmentResponse{
		ID:       establishment.ID,
		RUC:      establishment.RUC,
		Name:     establishment.Name,
		Phone:    establishment.Phone,
		Address:  establishment.Address,
		ImageUrl: establishment.ImageUrl,
		IsActive: establishment.IsActive,
		Admin:    adminResponse,
		AdminID:  establishment.AdminID,
	}

	return establishmentResponse, nil
}

// UpdateEstablishmentByAdminID updates the establishment associated with the admin.
func (s *establishmentService) UpdateEstablishmentByAdminID(adminID uint, req request.UpdateEstablishmentRequest) (*response.EstablishmentResponse, error) {
	establishment, err := s.establishmentRepo.GetEstablishmentByAdminID(adminID)
	if err != nil {
		return nil, err
	}

	// Update fields from the request
	establishment.RUC = req.RUC
	establishment.Name = req.Name
	establishment.Phone = req.Phone
	establishment.Address = req.Address
	establishment.ImageUrl = req.ImageUrl
	establishment.IsActive = req.IsActive
	establishment.LateFeePercentage = req.LateFeePercentage

	if err := s.establishmentRepo.UpdateEstablishment(establishment); err != nil {
		return nil, err
	}

	admin, err := s.userRepo.GetUserByID(adminID)

	if err != nil {
		return nil, err
	}

	adminResponse := &response.UserResponse{
		ID:      admin.ID,
		DNI:     admin.DNI,
		Email:   admin.Email,
		Name:    admin.Name,
		Rol:     admin.Rol,
		Address: admin.Address,
		Phone:   admin.Phone,
	}

	return establishmentToResponse(establishment, adminResponse), nil
}

// UploadEstablishmentLogo uploads an establishment logo and returns the URL.
func (s *establishmentService) UploadEstablishmentLogo(file *multipart.FileHeader) (string, error) {
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

	// 2. File Size Validation (Example: 2MB limit)
	if file.Size > 2*1024*1024 {
		return "", ErrFileSizeTooLarge
	}

	// 3. Create the "establishments_images" directory if it doesn't exist
	imagesDir := "establishments_images"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.Mkdir(imagesDir, 0755)
		if err != nil {
			return "", err
		}
	}

	// 4. Generate a unique filename (you can use UUIDs or a timestamp)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt) // Example using timestamp

	// 5. Create the full image file path
	imagePath := filepath.Join(imagesDir, newFilename)

	// 6. Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening uploaded file: %w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Println("error closing file:", err)
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
			fmt.Println("error closing file:", err)
		}
	}(dst)

	// 8. Copy the uploaded file contents to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("error copying image: %w", err)
	}

	// 9. Return the URL of the uploaded image
	return imagePath, nil
}
