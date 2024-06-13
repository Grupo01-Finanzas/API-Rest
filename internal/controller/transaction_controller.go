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

// TransactionController handles API requests related to transactions.
type TransactionController struct {
	transactionService service.TransactionService
}

// NewTransactionController creates a new instance of TransactionController.
func NewTransactionController(transactionService service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// CreateTransaction godoc
// @Summary Create Transaction
// @Description Create a new transaction for a credit account.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param transaction body request.CreateTransactionRequest true "Transaction Data"
// @Success 201 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/transactions [post]
func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req request.CreateTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := c.transactionService.CreateTransaction(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

// GetTransactionByID godoc
// @Summary Get Transaction by ID
// @Description Get a transaction by its ID.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/transactions/{id} [get]
func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid Transaction ID"})
		return
	}

	resp, err := c.transactionService.GetTransactionByID(uint(transactionID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Transaction not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// UpdateTransaction godoc
// @Summary Update Transaction
// @Description Update a transaction by its ID.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Param transaction body request.UpdateTransactionRequest true "Transaction Data"
// @Success 200 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/transactions/{id} [put]
func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid Transaction ID"})
		return
	}

	var req request.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := c.transactionService.UpdateTransaction(uint(transactionID), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Transaction not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// DeleteTransaction godoc
// @Summary Delete Transaction
// @Description Delete a transaction by its ID.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param id path int true "Transaction ID"
// @Success 204 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/transactions/{id} [delete]
func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid Transaction ID"})
		return
	}

	if err := c.transactionService.DeleteTransaction(uint(transactionID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Transaction not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// GetTransactionsByCreditAccountID godoc
// @Summary Get Transaction by Credit Account ID
// @Description Get all transactions for a specific credit account.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param creditAccountID path int true "Credit Account ID"
// @Success 200 {array} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/credit-accounts/{creditAccountID}/transactions [get]
func (c *TransactionController) GetTransactionsByCreditAccountID(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid Credit Account ID"})
		return
	}

	resp, err := c.transactionService.GetTransactionsByCreditAccountID(uint(creditAccountID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit Account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
