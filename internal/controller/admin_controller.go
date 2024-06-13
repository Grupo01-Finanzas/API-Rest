package controller

import (
	"ApiRestFinance/internal/model/entities"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService service.AdminService
}

func NewAdminController(adminService service.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

// CreateAdmin godoc
// @Summary      Create a new admin
// @Description  Creates a new admin with the provided data
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        admin  body      request.CreateAdminRequest  true  "Admin data"
// @Success      201     {object}  response.AdminResponse
// @Failure      400     {object}  map[string]string  "Invalid request"
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /api/v1/admins [post]
func (c *AdminController) CreateAdmin(ctx *gin.Context) {
	var req request.CreateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin := &entities.Admin{
		UserID:          req.UserID,
		EstablishmentID: req.EstablishmentID,
		IsActive:        req.IsActive,
	}

	if err := c.adminService.CreateAdmin(admin); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the created admin
	createdAdmin, err := c.adminService.GetAdminByUserID(req.UserID) // Call directly on adminRepo
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Admin created but failed to retrieve: " + err.Error()})
		return
	}

	// Get the associated establishment using the new method in adminService
	establishment, err := c.adminService.GetEstablishmentByID(createdAdmin.EstablishmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve associated establishment: " + err.Error()})
		return
	}

	resp := response.AdminResponse{
		ID:              createdAdmin.ID,
		UserID:          createdAdmin.UserID,
		EstablishmentID: createdAdmin.EstablishmentID,
		IsActive:        createdAdmin.IsActive,
		CreatedAt:       createdAdmin.CreatedAt,
		UpdatedAt:       createdAdmin.UpdatedAt,
		Establishment: response.EstablishmentResponse{ // Include establishment details
			ID:        establishment.ID,
			RUC:       establishment.RUC,
			Name:      establishment.Name,
			Phone:     establishment.Phone,
			Address:   establishment.Address,
			IsActive:  establishment.IsActive,
			CreatedAt: establishment.CreatedAt,
			UpdatedAt: establishment.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetAllAdmins godoc
// @Summary      Get all admins
// @Description  Gets a list of all admins.
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Success      200     {array}   response.AdminResponse
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /api/v1/admins [get]
func (c *AdminController) GetAllAdmins(ctx *gin.Context) {
	admins, err := c.adminService.GetAllAdmins()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []response.AdminResponse
	for _, admin := range admins {
		// Get the associated establishment
		establishment, err := c.adminService.GetEstablishmentByID(admin.EstablishmentID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve associated establishment: " + err.Error()})
			return
		}
		// Get the associated admin
		admin, err := c.adminService.GetAdminByUserID(admin.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve associated admin: " + err.Error()})
			return
		}

		resp = append(resp, response.AdminResponse{
			ID:              admin.ID,
			UserID:          admin.UserID,
			EstablishmentID: admin.EstablishmentID,
			IsActive:        admin.IsActive,
			CreatedAt:       admin.CreatedAt,
			UpdatedAt:       admin.UpdatedAt,
			Establishment: response.EstablishmentResponse{
				ID:        establishment.ID,
				RUC:       establishment.RUC,
				Name:      establishment.Name,
				Phone:     establishment.Phone,
				Address:   establishment.Address,
				IsActive:  establishment.IsActive,
				CreatedAt: establishment.CreatedAt,
				UpdatedAt: establishment.UpdatedAt,
			},
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetAdminByID godoc
// @Summary      Get admin by ID
// @Description  Gets an admin by their ID.
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Admin ID"
// @Success      200  {object}  response.AdminResponse
// @Failure      400  {object}  map[string]string  "Invalid admin ID"
// @Failure      404  {object}  map[string]string  "Admin not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /api/v1/admins/{id} [get]
func (c *AdminController) GetAdminByID(ctx *gin.Context) {
	adminID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	admin, err := c.adminService.GetAdminByID(uint(adminID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Get the associated establishment
	establishment, err := c.adminService.GetEstablishmentByID(admin.EstablishmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve associated establishment: " + err.Error()})
		return
	}

	resp := response.AdminResponse{
		ID:              admin.ID,
		UserID:          admin.UserID,
		EstablishmentID: admin.EstablishmentID,
		IsActive:        admin.IsActive,
		CreatedAt:       admin.CreatedAt,
		UpdatedAt:       admin.UpdatedAt,
		Establishment: response.EstablishmentResponse{
			ID:        establishment.ID,
			RUC:       establishment.RUC,
			Name:      establishment.Name,
			Phone:     establishment.Phone,
			Address:   establishment.Address,
			IsActive:  establishment.IsActive,
			CreatedAt: establishment.CreatedAt,
			UpdatedAt: establishment.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateAdmin godoc
// @Summary      Update admin
// @Description  Updates an admin\'s data.
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        id     path      int                     true  "Admin ID"
// @Param        admin  body      request.UpdateAdminRequest  true  "Updated admin data"
// @Success      200     {object}  map[string]string  "Admin updated successfully"
// @Failure      400     {object}  map[string]string  "Invalid admin ID or request body"
// @Failure      404     {object}  map[string]string  "Admin not found"
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /api/v1/admins/{id} [put]
func (c *AdminController) UpdateAdmin(ctx *gin.Context) {
	adminID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	var req request.UpdateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the existing admin
	admin, err := c.adminService.GetAdminByID(uint(adminID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Update the admin fields (only IsActive in this case)
	admin.IsActive = req.IsActive

	// Save the updated admin
	if err := c.adminService.UpdateAdmin(admin); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})
}

// DeleteAdmin godoc
// @Summary      Delete admin
// @Description  Deletes an admin by their ID.
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Admin ID"
// @Success      200  {object}  map[string]string  "Admin deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid admin ID"
// @Failure      404  {object}  map[string]string  "Admin not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /api/v1/admins/{id} [delete]
func (c *AdminController) DeleteAdmin(ctx *gin.Context) {
	adminID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin ID"})
		return
	}

	if err := c.adminService.DeleteAdmin(uint(adminID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}

// RegisterEstablishment godoc
// @Summary      Register establishment
// @Description  Registers a new establishment associated with the admin.
// @Tags         Admins
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                          true  "Bearer {accessToken}"
// @Param        establishment  body      request.CreateEstablishmentRequest  true  "Establishment data"
// @Success      201     {object}  response.EstablishmentResponse
// @Failure      400     {object}  map[string]string  "Invalid request"
// @Failure      401     {object}  map[string]string  "Unauthorized"
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /api/v1/establishments [post]
func (c *AdminController) RegisterEstablishment(ctx *gin.Context) {
	var req request.CreateEstablishmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userIDFloat, ok := claimsMap["user_id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID := uint(userIDFloat)

	establishment := &entities.Establishment{
		RUC:      req.RUC,
		Name:     req.Name,
		Phone:    req.Phone,
		Address:  req.Address,
		IsActive: req.IsActive,
	}

	if err := c.adminService.RegisterEstablishment(establishment, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the created establishment by ID (you might need to update your establishmentRepo)
	createdEstablishment, err := c.adminService.GetEstablishmentByID(establishment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Establishment created but failed to retrieve: " + err.Error()})
		return
	}

	resp := response.EstablishmentResponse{
		ID:        createdEstablishment.ID,
		RUC:       createdEstablishment.RUC,
		Name:      createdEstablishment.Name,
		Phone:     createdEstablishment.Phone,
		Address:   createdEstablishment.Address,
		IsActive:  createdEstablishment.IsActive,
		CreatedAt: createdEstablishment.CreatedAt,
		UpdatedAt: createdEstablishment.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, resp)
}
