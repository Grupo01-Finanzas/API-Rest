package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UserService handles user-related operations.
type UserService interface {
	CreateClient(req request.CreateClientRequest) (*response.UserResponse, error)
	GetUserByID(userID uint) (*response.UserResponse, error)
	UpdateUser(userID uint, req request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteUser(userID uint) error
	GetClientsByEstablishmentID(establishmentID uint) ([]entities.User, error)
	UploadUserPhoto(photo *multipart.FileHeader, userID uint) (string, error)
}

type userService struct {
	userRepo          repository.UserRepository
	creditAccountRepo repository.CreditAccountRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repository.UserRepository, creditAccountRepo repository.CreditAccountRepository) UserService {
	return &userService{userRepo: userRepo, creditAccountRepo: creditAccountRepo}
}

// CreateClient creates a new client user and their associated credit account.
func (s *userService) CreateClient(req request.CreateClientRequest) (*response.UserResponse, error) {
	// Create the User entity
	user := &entities.User{
		DNI:      req.DNI,
		Email:    req.Email,
		Password: generateRandomPassword(),
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		Rol:      enums.CLIENT,
	}

	// Create the CreditAccount entity
	creditAccount := &entities.CreditAccount{
		EstablishmentID:         req.EstablishmentID,
		ClientID:                user.ID,
		CreditLimit:             req.CreditLimit,
		MonthlyDueDate:          req.MonthlyDueDate,
		InterestRate:            req.InterestRate,
		InterestType:            req.InterestType,
		CreditType:              req.CreditType,
		GracePeriod:             req.GracePeriod,
		IsBlocked:               false,
		LastInterestAccrualDate: time.Now(),
		CurrentBalance:          0.0,
		LateFeePercentage:       req.LateFeePercentage,
	}

	// Use the CreditAccountRepository to handle the creation in a transaction
	if err := s.creditAccountRepo.CreateClientAndCreditAccount(user, creditAccount); err != nil {
		return nil, fmt.Errorf("error during client creation: %w", err)
	}

	return _NewUserResponse(user), nil
}

// GetUserByID retrieves a user by their ID.
func (s *userService) GetUserByID(userID uint) (*response.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return _NewUserResponse(user), nil
}

// UpdateUser updates an existing user.
func (s *userService) UpdateUser(userID uint, req request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Update user fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	// Update the PhotoUrl if provided
	if req.PhotoUrl != "" {
		user.PhotoUrl = req.PhotoUrl
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return NewUserResponse(user), nil
}

// DeleteUser deletes a user and their associated credit account.
func (s *userService) DeleteUser(userID uint) error {
	// You might want to add checks here to ensure you are deleting the correct type of user (CLIENT)
	return s.creditAccountRepo.DeleteClientAndCreditAccount(userID)
}

// GetClientsByEstablishmentID retrieves all users with the CLIENT role
// associated with a specific establishment.
func (s *userService) GetClientsByEstablishmentID(establishmentID uint) ([]entities.User, error) {
	return s.userRepo.GetClientsByEstablishmentID(establishmentID)
}

// UploadUserPhoto handles the actual photo upload to the server.
func (s *userService) UploadUserPhoto(photo *multipart.FileHeader, userID uint) (string, error) {
	// 1. File Type Validation (Only allow images)
	allowedFileTypes := []string{".jpg", ".jpeg", ".png", ".gif"}

	fileExt := strings.ToLower(filepath.Ext(photo.Filename))
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
	if photo.Size > 2*1024*1024 {
		return "", ErrFileSizeTooLarge
	}

	// 3. Create the images directory if it doesn't exist
	imagesDir := "images_user"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.Mkdir(imagesDir, 0755)
		if err != nil {
			return "", err
		}
	}

	// 4. Generate a unique filename for the image (you can use UUIDs or any other method)
	newFilename := fmt.Sprintf("%d%s", userID, fileExt)

	// 5. Create the full path to the image file
	imagePath := filepath.Join(imagesDir, newFilename)

	// 6. Open the uploaded photo file
	file, err := photo.Open()
	if err != nil {
		return "", fmt.Errorf("error opening photo file: %w", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file:", err)
		}
	}(file)

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
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error copying photo: %w", err)
	}

	// 9. Return the relative URL of the uploaded image
	return imagePath, nil
}

// NewUserResponse converts a User entity to a UserResponse DTO.
func _NewUserResponse(user *entities.User) *response.UserResponse {
	if user == nil {
		return nil
	}
	return &response.UserResponse{
		ID:        user.ID,
		DNI:       user.DNI,
		Email:     user.Email,
		Name:      user.Name,
		Address:   user.Address,
		Phone:     user.Phone,
		PhotoUrl:  user.PhotoUrl,
		Rol:       user.Rol,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
