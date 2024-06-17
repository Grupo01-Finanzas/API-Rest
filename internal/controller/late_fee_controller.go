package controller

import (
	"errors"
	"net/http"
	"strconv"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LateFeeController handles API requests related to late fees.
type LateFeeController struct {
	lateFeeService service.LateFeeService
}

// NewLateFeeController creates a new LateFeeController instance.
func NewLateFeeController(lateFeeService service.LateFeeService) *LateFeeController {
	return &LateFeeController{lateFeeService: lateFeeService}
}

// CreateLateFee godoc
// @Summary Create a late fee
// @Description Creates a new late fee for a credit account.
// @Tags LateFees
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param lateFee body request.CreateLateFeeRequest true "Late fee details"
// @Success 201 {object} response.LateFeeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fees [post]
func (c *LateFeeController) CreateLateFee(ctx *gin.Context) {
	var req request.CreateLateFeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	lateFee, err := c.lateFeeService.CreateLateFee(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, lateFee)
}

// GetLateFeeByID godoc
// @Summary Get late fee by ID
// @Description Retrieves a late fee by its ID.
// @Tags LateFees
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Late Fee ID"
// @Success 200 {object} response.LateFeeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fees/{id} [get]
func (c *LateFeeController) GetLateFeeByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid late fee ID"})
		return
	}

	lateFee, err := c.lateFeeService.GetLateFeeByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Late fee not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFee)
}

// UpdateLateFee godoc
// @Summary Update a late fee
// @Description Updates an existing late fee by its ID.
// @Tags LateFees
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Late Fee ID"
// @Param lateFee body request.UpdateLateFeeRequest true "Updated late fee details"
// @Success 200 {object} response.LateFeeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fees/{id} [put]
func (c *LateFeeController) UpdateLateFee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid late fee ID"})
		return
	}

	var req request.UpdateLateFeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	lateFee, err := c.lateFeeService.UpdateLateFee(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Late fee not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFee)
}

// DeleteLateFee godoc
// @Summary Delete a late fee
// @Description Deletes a late fee by its ID.
// @Tags LateFees
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Late Fee ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /late-fees/{id} [delete]
func (c *LateFeeController) DeleteLateFee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid late fee ID"})
		return
	}

	if err := c.lateFeeService.DeleteLateFee(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Late fee not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetLateFeesByCreditAccountID godoc
// @Summary Get late fees by credit account ID
// @Description Retrieves all late fees for a specific credit account.
// @Tags LateFees
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Success 200 {array} response.LateFeeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/late-fees [get]
func (c *LateFeeController) GetLateFeesByCreditAccountID(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	lateFees, err := c.lateFeeService.GetLateFeesByCreditAccountID(uint(creditAccountID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lateFees)
}
