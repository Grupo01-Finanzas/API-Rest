package controller

import (
	"ApiRestFinance/internal/model/entities"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"

	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/service"

	"github.com/gin-gonic/gin"
)

// UserController handles all user-related endpoints, including Admins and Clients
type UserController struct {
	userService          service.UserService
	adminService         service.AdminService
	creditAccountService service.CreditAccountService
	establishmentService service.EstablishmentService
}

// NewUserController creates a new instance of UserController.
func NewUserController(userService service.UserService, adminService service.AdminService, creditAccountService service.CreditAccountService, establishmentService service.EstablishmentService) *UserController {
	return &UserController{userService: userService, adminService: adminService, creditAccountService: creditAccountService, establishmentService: establishmentService}
}

// CreateClient godoc
// @Summary      Create Client
// @Description  Creates a new client user with an associated credit account. Only Admins can create clients.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                  true  "Bearer {token}"
// @Param        client         body      request.CreateClientRequest  true  "Client data"
// @Success      201  {object}  response.UserResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients [post]
func (c *UserController) CreateClient(ctx *gin.Context) {
	var req request.CreateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can create clients"})
		return
	}

	userId := middleware.GetUserIDFromContext(ctx)

	establishment, err := c.establishmentService.GetEstablishmentByAdminID(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}
	req.EstablishmentID = establishment.ID

	userResponse, err := c.userService.CreateClient(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, userResponse)
}

// UpdatePassword godoc
// @Summary      Update Client Password
// @Description  Updates the password for the authenticated client.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                      true  "Bearer {token}"
// @Param        newPassword     body      request.ResetPasswordRequest  true  "New password data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Router       /clients/me/password [put]
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var req request.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the authenticated user is a CLIENT
	if middleware.GetUserRoleFromContext(ctx) != enums.CLIENT {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only clients can update their password"})
		return
	}

	userID := middleware.GetUserIDFromContext(ctx)

	err := c.userService.UpdatePassword(userID, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// GetUserByID godoc
// @Summary      Get User by ID
// @Description  Retrieves a user by their ID. Admins can retrieve any user, Clients can only retrieve themselves.
// @Tags         Users
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id             path      int  true  "User ID"
// @Success      200  {object}  response.UserResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	authUserID := middleware.GetUserIDFromContext(ctx)
	authUserRole := middleware.GetUserRoleFromContext(ctx)

	// Authorization: Admins can access any user; Clients can only access their own data
	if authUserRole != enums.ADMIN && authUserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Not authorized to access this user"})
		return
	}

	userResponse, err := c.userService.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// DeleteUser godoc
// @Summary      Delete User
// @Description  Deletes a user by their ID. Only Admins can delete users.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id             path      int  true  "User ID"
// @Success      204  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can delete users"})
		return
	}

	if err := c.userService.DeleteUser(uint(userID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "User deleted successfully"})
}

// GetAdminProfile godoc
// @Summary      Get Admin Profile
// @Description  Retrieves the profile information of the authenticated admin.
// @Tags         Users
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  response.AdminResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admins/me [get]
func (c *UserController) GetAdminProfile(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can access this endpoint"})
		return
	}

	adminResponse, err := c.adminService.GetAdminByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, adminResponse)
}

// UpdateAdminProfile godoc
// @Summary      Update Admin Profile
// @Description  Updates the profile information of the authenticated admin.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        admin          body      request.UpdateUserRequest  true  "Updated admin data"
// @Success      200  {object}  response.AdminResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admins/me [put]
func (c *UserController) UpdateAdminProfile(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can access this endpoint"})
		return
	}

	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	adminResponse, err := c.adminService.UpdateAdmin(userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, adminResponse)
}

// UpdateClientCreditAccount godoc
// @Summary      Update Client Credit Account
// @Description  Updates credit account details of a client. Only Admins can update credit accounts.
// @Tags         Credit Accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                        true  "Bearer {token}"
// @Param        clientID       path      int                        true  "Client User ID"
// @Param        creditAccount  body      request.UpdateCreditAccountRequest  true  "Updated credit account data"
// @Success      200  {object}  response.CreditAccountResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/{clientID}/credit-account [put]
func (c *UserController) UpdateClientCreditAccount(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("clientID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req request.UpdateCreditAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can update credit accounts"})
		return
	}

	creditAccountResponse, err := c.creditAccountService.UpdateCreditAccountByClientID(uint(clientID), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccountResponse)
}

// GetClientsByEstablishmentID godoc
// @Summary      Get Clients by Establishment ID
// @Description  Gets all clients associated with an establishment. Only Admins can access this endpoint.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        establishmentID   path      int  true  "Establishment ID"
// @Success      200  {array}   response.UserResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /establishments/{establishmentID}/clients [get]
func (c *UserController) GetClientsByEstablishmentID(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	// Ensure the authenticated user is an ADMIN
	authUserRole := middleware.GetUserRoleFromContext(ctx)
	if authUserRole != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can access clients"})
		return
	}

	clients, err := c.userService.GetClientsByEstablishmentID(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Convert to UserResponse
	var userResponses []response.UserResponse
	for _, client := range clients {
		userResponses = append(userResponses, *_NewUserResponse(&client))
	}

	ctx.JSON(http.StatusOK, userResponses)
}

// UploadUserPhoto godoc
// @Summary      Upload User PhotoUrl
// @Description  Uploads a profile photo for a user.
// @Tags         Users
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id             path      int                      true  "User ID"
// @Param        photo          formData      file  true  "User profile photo"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users/{id}/photo [post]
func (c *UserController) UploadUserPhoto(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Handle file upload
	file, err := ctx.FormFile("photo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Error uploading file: " + err.Error()})
		return
	}

	// Authorization Logic:
	authUserID := middleware.GetUserIDFromContext(ctx)
	authUserRole := middleware.GetUserRoleFromContext(ctx)

	// Allow a user to update their own photo or an admin to update any user's photo
	if authUserRole != enums.ADMIN && authUserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Forbidden: You are not authorized to upload a photo for this user"})
		return
	}

	// Upload photo using the service
	photoURL, err := c.userService.UploadUserPhoto(file, uint(userID))
	if err != nil {
		// Handle errors (file type, size, storage errors)
		if errors.Is(err, service.ErrInvalidFileType) || errors.Is(err, service.ErrFileSizeTooLarge) {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error uploading photo: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"photo_url": photoURL})
}

// UpdateUser godoc
// @Summary      Update User
// @Description  Updates user details, including the URL of the profile photo.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                  true  "Bearer {token}"
// @Param        id             path      int                      true  "User ID"
// @Param        user           body      request.UpdateUserRequest  true  "Updated user data (including photo_url)"
// @Success      200  {object}  response.UserResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Authorization Logic:
	authUserID := middleware.GetUserIDFromContext(ctx)
	authUserRole := middleware.GetUserRoleFromContext(ctx)

	// Allow admins to update any user, but clients can only update themselves
	if authUserRole != enums.ADMIN && authUserID != uint(userID) {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Forbidden: You are not authorized to update this user"})
		return
	}

	// Update user using the service (which handles photo uploads)
	userResponse, err := c.userService.UpdateUser(uint(userID), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error updating user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

// GetUserIDByEmail godoc
// @Summary      Get User ID by Email
// @Description  Retrieves the ID of a user by their email address. This endpoint is typically for internal use or admin purposes.
// @Tags         Users
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        email          query       string  true  "User's email address"
// @Success      200  {object}  map[string]uint
// @Failure      400  {object}  response.ErrorResponse  "Invalid email format"
// @Failure      401  {object}  response.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse  "Forbidden (only for admins)"
// @Failure      404  {object}  response.ErrorResponse  "User not found"
// @Router       /users/email-to-id [get]
func (c *UserController) GetUserIDByEmail(ctx *gin.Context) {
	email := ctx.Query("email")

	// Input Validation
	if !isValidEmailFormat(email) {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid email format"})
		return
	}

	// Authorization (Optional - If you want to restrict access to admins only)
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Forbidden: Only admins can access this endpoint"})
		return
	}

	userID, err := c.userService.GetUserIDByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

// Helper Function for Email Validation
func isValidEmailFormat(email string) bool {
	// You can use a more robust email validation library here if needed
	return strings.Contains(email, "@")
}

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
