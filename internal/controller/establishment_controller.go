package controller

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"

	"ApiRestFinance/internal/service"
)

type EstablishmentController struct {
	establishmentService service.EstablishmentService
}

func NewEstablishmentController(establishmentService service.EstablishmentService) *EstablishmentController {
	return &EstablishmentController{
		establishmentService: establishmentService,
	}
}

// GetAllEstablishments godoc
// @Summary Get all establishments
// @Description Retrieve a list of all establishments.
// @Tags Establishments
// @Produce json
// @Success 200 {array} response.EstablishmentResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments [get]
func (c *EstablishmentController) GetAllEstablishments(ctx *gin.Context) {
	establishments, err := c.establishmentService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, establishments)
}

// GetEstablishmentByID godoc
// @Summary Get an establishment by ID
// @Description Retrieve an establishment by its ID.
// @Tags Establishments
// @Produce json
// @Param id path uint true "Establishment ID"
// @Success 200 {object} response.EstablishmentResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments/{id} [get]
func (c *EstablishmentController) GetEstablishmentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	establishment, err := c.establishmentService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Establishment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, establishment)
}

// CreateEstablishment godoc
// @Summary Create a new establishment
// @Description Create a new establishment with admin details.
// @Tags Establishments
// @Accept json
// @Produce json
// @Param establishment body request.CreateEstablishmentRequest true "Establishment details"
// @Success 201 {object} response.EstablishmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments [post]
func (c *EstablishmentController) CreateEstablishment(ctx *gin.Context) {
	var req request.CreateEstablishmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	establishment, err := c.establishmentService.Create(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, establishment)
}

// UpdateEstablishment godoc
// @Summary Update an existing establishment
// @Description Update an existing establishment by ID.
// @Tags Establishments
// @Accept json
// @Produce json
// @Param id path uint true "Establishment ID"
// @Param establishment body request.UpdateEstablishmentRequest true "Establishment details"
// @Success 200 {object} response.EstablishmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments/{id} [put]
func (c *EstablishmentController) UpdateEstablishment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	var req request.UpdateEstablishmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	establishment, err := c.establishmentService.Update(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Establishment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, establishment)
}

// DeleteEstablishment godoc
// @Summary Delete an establishment
// @Description Delete an establishment by ID.
// @Tags Establishments
// @Produce json
// @Param id path uint true "Establishment ID"
// @Success 204
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments/{id} [delete]
func (c *EstablishmentController) DeleteEstablishment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	err = c.establishmentService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Establishment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// RegisterProducts godoc
// @Summary Register products for an establishment
// @Description Associate a list of products with an establishment.
// @Tags Establishments
// @Accept json
// @Produce json
// @Param id path uint true "Establishment ID"
// @Param products body []uint true "List of product IDs"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments/{id}/products [post]
func (c *EstablishmentController) RegisterProducts(ctx *gin.Context) {
	establishmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	var productIDs []uint
	if err := ctx.ShouldBindJSON(&productIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	err = c.establishmentService.RegisterProducts(uint(establishmentID), productIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Establishment or product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// AddClientToEstablishment godoc
// @Summary Add a client to an establishment
// @Description Associate a client with an establishment.
// @Tags Establishments
// @Produce json
// @Param id path uint true "Establishment ID"
// @Param client_id path uint true "Client ID"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/establishments/{id}/clients/{client_id} [put]
func (c *EstablishmentController) AddClientToEstablishment(ctx *gin.Context) {
	establishmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	clientID, err := strconv.ParseUint(ctx.Param("client_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	err = c.establishmentService.AddClientToEstablishment(uint(establishmentID), uint(clientID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Establishment or client not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
