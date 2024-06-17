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

// InstallmentController handles API requests related to installments.
type InstallmentController struct {
	installmentService service.InstallmentService
}

// NewInstallmentController creates a new InstallmentController.
func NewInstallmentController(installmentService service.InstallmentService) *InstallmentController {
	return &InstallmentController{installmentService: installmentService}
}

// CreateInstallment godoc
// @Summary Create a new installment
// @Description Creates a new installment for a credit account.
// @Tags Installments
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param installment body request.CreateInstallmentRequest true "Installment details"
// @Success 201 {object} response.InstallmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /installments [post]
func (c *InstallmentController) CreateInstallment(ctx *gin.Context) {
	var req request.CreateInstallmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	installment, err := c.installmentService.CreateInstallment(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, installment)
}

// GetInstallmentByID godoc
// @Summary Get installment by ID
// @Description Retrieves an installment by its ID.
// @Tags Installments
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Installment ID"
// @Success 200 {object} response.InstallmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /installments/{id} [get]
func (c *InstallmentController) GetInstallmentByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid installment ID"})
		return
	}

	installment, err := c.installmentService.GetInstallmentByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Installment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installment)
}

// UpdateInstallment godoc
// @Summary Update an installment
// @Description Updates an existing installment by its ID.
// @Tags Installments
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Installment ID"
// @Param installment body request.UpdateInstallmentRequest true "Updated installment details"
// @Success 200 {object} response.InstallmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /installments/{id} [put]
func (c *InstallmentController) UpdateInstallment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid installment ID"})
		return
	}

	var req request.UpdateInstallmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	installment, err := c.installmentService.UpdateInstallment(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Installment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installment)
}

// DeleteInstallment godoc
// @Summary Delete an installment
// @Description Deletes an installment by its ID.
// @Tags Installments
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Installment ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /installments/{id} [delete]
func (c *InstallmentController) DeleteInstallment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid installment ID"})
		return
	}

	if err := c.installmentService.DeleteInstallment(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Installment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetInstallmentsByCreditAccountID godoc
// @Summary Get installments by credit account ID
// @Description Retrieves all installments for a specific credit account.
// @Tags Installments
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Success 200 {array} response.InstallmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/installments [get]
func (c *InstallmentController) GetInstallmentsByCreditAccountID(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	installments, err := c.installmentService.GetInstallmentsByCreditAccountID(uint(creditAccountID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installments)
}

// GetOverdueInstallments godoc
// @Summary Get overdue installments by credit account ID
// @Description Retrieves all overdue installments for a specific credit account.
// @Tags Installments
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Success 200 {array} response.InstallmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/installments/overdue [get]
func (c *InstallmentController) GetOverdueInstallments(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	installments, err := c.installmentService.GetOverdueInstallments(uint(creditAccountID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installments)
}
