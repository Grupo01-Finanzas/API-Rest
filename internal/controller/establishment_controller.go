package controller

import (
	"net/http"
	"strconv"

	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"

	"github.com/gin-gonic/gin"
)

// EstablishmentController handles establishment-related endpoints.
type EstablishmentController struct {
	establishmentService service.EstablishmentService
}

// NewEstablishmentController creates a new instance of EstablishmentController.
func NewEstablishmentController(establishmentService service.EstablishmentService) *EstablishmentController {
	return &EstablishmentController{establishmentService: establishmentService}
}

// CreateEstablishment godoc
// @Summary      Create Establishment
// @Description  Creates a new establishment for the authenticated admin.
// @Tags         Establishments
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                          true  "Bearer {token}"
// @Param        establishment  body      request.CreateEstablishmentRequest  true  "Establishment data"
// @Success      201  {object}  response.EstablishmentResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /establishments [post]
func (c *EstablishmentController) CreateEstablishment(ctx *gin.Context) {
	var req request.CreateEstablishmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	adminID := middleware.GetUserIDFromContext(ctx) // Get authenticated admin\'s ID

	establishment, err := c.establishmentService.CreateEstablishment(&req, adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, establishment)
}

// GetEstablishment godoc
// @Summary      Get Establishment
// @Description  Gets the establishment details for the authenticated admin.
// @Tags         Establishments
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  response.EstablishmentResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /establishments/me [get]
func (c *EstablishmentController) GetEstablishment(ctx *gin.Context) {
	adminID := middleware.GetUserIDFromContext(ctx)

	establishment, err := c.establishmentService.GetEstablishmentByAdminID(adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, establishment)
}

// UpdateEstablishment godoc
// @Summary      Update Establishment
// @Description  Updates the establishment details for the authenticated admin.
// @Tags         Establishments
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string                          true  "Bearer {token}"
// @Param        establishment  body      request.UpdateEstablishmentRequest  true  "Updated establishment data"
// @Success      200  {object}  response.EstablishmentResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /establishments/me [put]
func (c *EstablishmentController) UpdateEstablishment(ctx *gin.Context) {
	var req request.UpdateEstablishmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	adminID := middleware.GetUserIDFromContext(ctx)

	establishment, err := c.establishmentService.UpdateEstablishmentByAdminID(adminID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, establishment)
}

// GetEstablishmentByID godoc
// @Summary      Get Establishment by ID
// @Description  Gets an establishment by its ID.
// @Tags         Establishments
// @Produce      json
// @Param        Authorization  header      string                          true  "Bearer {token}"
// @Param        establishmentID   path      int  true  "Establishment ID"
// @Success      200  {object}  response.EstablishmentResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /establishments/{establishmentID} [get]
func (c *EstablishmentController) GetEstablishmentByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	adminID := middleware.GetUserIDFromContext(ctx)
	establishment, err := c.establishmentService.GetEstablishmentByAdminID(adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	if establishment.ID != uint(id) {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Establishment not found"})
		return
	}

	ctx.JSON(http.StatusOK, establishment)
}
