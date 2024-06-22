package controller

import (
	"errors"
	"net/http"
	"strconv"

	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities/enums"
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
// @Summary      Create Transaction
// @Description  Create a new transaction (purchase or payment).
// @Tags         Transactions
// @Accept  json
// @Produce  json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param transaction body request.CreateTransactionRequest true "Transaction Data"
// @Success 201 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions [post]
func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req request.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate transaction type
	if req.TransactionType != enums.Purchase && req.TransactionType != enums.Payment {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid transaction type"})
		return
	}

	// Check user role
	userRole := middleware.GetUserRoleFromContext(ctx)

	if req.TransactionType == enums.Purchase && userRole != enums.CLIENT {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only clients can create purchase transactions"})
		return
	} else if req.TransactionType == enums.Payment && userRole != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can create payment transactions"})
		return
	}

	resp, err := c.transactionService.CreateTransaction(req)
	if err != nil {
		if errors.Is(err, service.ErrCreditAccountNotFound) || errors.Is(err, service.ErrInvalidTransactionType) || errors.Is(err, service.ErrInsufficientBalance) {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		}
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
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions/{id} [get]
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

	// Authorization: Only the admin or the client associated with the transaction can access it
	authUserID := middleware.GetUserIDFromContext(ctx)
	authUserRole := middleware.GetUserRoleFromContext(ctx)
	if authUserRole != enums.ADMIN && resp.CreditAccountID != authUserID { // Assuming CreditAccountID is the Client User ID
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Forbidden: Not authorized to access this transaction"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetTransactionsByCreditAccountID godoc
// @Summary Get Transaction by Credit Account ID
// @Description Get all transactions for a specific credit account.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Success 200 {array} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/transactions [get]
func (c *TransactionController) GetTransactionsByCreditAccountID(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid Credit Account ID"})
		return
	}

	// Authorization: Only the admin or the client associated with the credit account can access its transactions
	authUserID := middleware.GetUserIDFromContext(ctx)
	authUserRole := middleware.GetUserRoleFromContext(ctx)
	if authUserRole != enums.ADMIN && uint(creditAccountID) != authUserID {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Forbidden: Not authorized to access transactions for this credit account"})
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

// UpdateTransaction godoc
// @Summary Update Transaction
// @Description Update a transaction by its ID. Only admins can update transactions.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Transaction ID"
// @Param transaction body request.UpdateTransactionRequest true "Transaction Data"
// @Success 200 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions/{id} [put]
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

	// Only admins can update transactions
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can update transactions"})
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
// @Description Delete a transaction by its ID. Only admins can delete transactions.
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Transaction ID"
// @Success 204 {object} response.TransactionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /transactions/{id} [delete]
func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid Transaction ID"})
		return
	}

	// Only admins can delete transactions
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can delete transactions"})
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

// ConfirmPayment godoc
// @Summary      Confirm Payment
// @Description  Confirms a pending payment using a confirmation code. Only admins can confirm payments.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id              path      int  true  "Transaction ID"
// @Param        confirmation   body      map[string]string  true  "Confirmation code"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /transactions/{id}/confirm [post]
func (c *TransactionController) ConfirmPayment(ctx *gin.Context) {
	transactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid transaction ID"})
		return
	}

	var req map[string]string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid request format"})
		return
	}

	confirmationCode, ok := req["confirmation_code"]
	if !ok || confirmationCode == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Confirmation code is required"})
		return
	}

	// Only admins can confirm payments
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can confirm payments"})
		return
	}

	if err := c.transactionService.ConfirmPayment(uint(transactionID), confirmationCode); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment confirmed successfully"})
}
