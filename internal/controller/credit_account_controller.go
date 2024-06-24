package controller

import (
	"errors"
	"fmt"
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

// CreditAccountController handles endpoints related to credit accounts.
type CreditAccountController struct {
	creditAccountService service.CreditAccountService
	establishmentService service.EstablishmentService
}

// NewCreditAccountController creates a new instance of CreditAccountController.
func NewCreditAccountController(creditAccountService service.CreditAccountService, establishmentService service.EstablishmentService) *CreditAccountController {
	return &CreditAccountController{creditAccountService: creditAccountService, establishmentService: establishmentService}
}

// CreateCreditAccount godoc
// @Summary      Create Credit Account
// @Description  Creates a new credit account for a client.
// @Tags         Credit Accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        creditAccount  body      request.CreateCreditAccountRequest  true  "Credit account data"
// @Success      201  {object}  response.CreditAccountResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts [post]
func (c *CreditAccountController) CreateCreditAccount(ctx *gin.Context) {
	var req request.CreateCreditAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Only admins can create credit accounts
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can create credit accounts"})
		return
	}

	userId := middleware.GetUserIDFromContext(ctx)

	establishment, err := c.establishmentService.GetEstablishmentByAdminID(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}

	creditAccount, err := c.creditAccountService.CreateCreditAccount(req, establishment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, creditAccount)
}

// GetCreditAccountByID godoc
// @Summary      Get Credit Account by ID
// @Description  Gets a credit account by its ID.
// @Tags         Credit Accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Credit Account ID"
// @Success      200  {object}  response.CreditAccountResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/{id} [get]
func (c *CreditAccountController) GetCreditAccountByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	creditAccount, err := c.creditAccountService.GetCreditAccountByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccount)
}

// GetCreditAccountByClientID godoc
// @Summary      Get Credit Account by Client ID
// @Description  Retrieves a credit account associated with a specific client.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        clientID path int true "Client ID"
// @Success      200 {object}  response.CreditAccountResponse
// @Failure      400 {object}  response.ErrorResponse
// @Failure      404 {object}  response.ErrorResponse
// @Failure      500 {object}  response.ErrorResponse
// @Router       /clients/{clientID}/credit-account [get]
func (c *CreditAccountController) GetCreditAccountByClientID(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("clientID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	// Authorization: Admins can access any client's credit account, Clients can only access their own
	authUserID := middleware.GetUserIDFromContext(ctx)
	authUserRole := middleware.GetUserRoleFromContext(ctx)
	if authUserRole != enums.ADMIN && authUserID != uint(clientID) {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Not authorized to access this credit account"})
		return
	}

	creditAccount, err := c.creditAccountService.GetCreditAccountByClientID(uint(clientID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found for this client"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccount)
}

// UpdateCreditAccount godoc
// @Summary      Update Credit Account
// @Description  Updates a credit account by its ID. Only Admins can update credit accounts.
// @Tags         Credit Accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id     path      int                      true  "Credit Account ID"
// @Param        creditAccount  body      request.UpdateCreditAccountRequest  true  "Updated credit account data"
// @Success      200     {object}  response.CreditAccountResponse
// @Failure      400     {object}  response.ErrorResponse
// @Failure      401     {object}  response.ErrorResponse
// @Failure      403     {object}  response.ErrorResponse
// @Failure      404     {object}  response.ErrorResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /credit-accounts/{id} [put]
func (c *CreditAccountController) UpdateCreditAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
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

	creditAccount, err := c.creditAccountService.UpdateCreditAccount(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccount)
}

// DeleteCreditAccount godoc
// @Summary      Delete Credit Account
// @Description  Deletes a credit account by its ID. Only Admins can delete credit accounts.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Credit Account ID"
// @Success      204  "No Content"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /credit-accounts/{id} [delete]
func (c *CreditAccountController) DeleteCreditAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	// Only Admins can delete credit accounts
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can delete credit accounts"})
		return
	}

	if err := c.creditAccountService.DeleteCreditAccount(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetCreditAccountsByEstablishmentID godoc
// @Summary      Get Credit Accounts by Establishment ID
// @Description  Retrieves all credit accounts associated with an establishment. Only Admins can access this endpoint.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        establishmentID path int true "Establishment ID"
// @Success      200 {array} response.CreditAccountResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /establishments/{establishmentID}/credit-accounts [get]
func (c *CreditAccountController) GetCreditAccountsByEstablishmentID(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	authUserRole := middleware.GetUserRoleFromContext(ctx)
	if authUserRole != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can access credit accounts"})
		return
	}

	creditAccounts, err := c.creditAccountService.GetCreditAccountsByEstablishmentID(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccounts)
}

// ApplyInterestToAccount godoc
// @Summary      Apply Interest to Account
// @Description  Applies interest to a specific credit account. Only Admins can apply interest.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        creditAccountID path int true "Credit Account ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/{creditAccountID}/apply-interest [post]
func (c *CreditAccountController) ApplyInterestToAccount(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	// Only Admins can apply interest
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can apply interest to credit accounts"})
		return
	}

	if err := c.creditAccountService.ApplyInterestToAccount(uint(creditAccountID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interest applied successfully"})
}

// ApplyLateFeeToAccount godoc
// @Summary      Apply Late Fee to Account
// @Description  Applies a late fee to a specific credit account. Only Admins can apply late fees.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        creditAccountID path int true "Credit Account ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/{creditAccountID}/apply-late-fee [post]
func (c *CreditAccountController) ApplyLateFeeToAccount(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	// Only Admins can apply late fees
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can apply late fees to credit accounts"})
		return
	}

	if err := c.creditAccountService.ApplyLateFeeToAccount(uint(creditAccountID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Late fee applied successfully"})
}

// GetOverdueCreditAccounts godoc
// @Summary      Get Overdue Credit Accounts
// @Description  Retrieves all overdue credit accounts for the authenticated admin\'s establishment.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {array}   response.CreditAccountResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/overdue [get]
func (c *CreditAccountController) GetOverdueCreditAccounts(ctx *gin.Context) {
	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can access this endpoint"})
		return
	}

	userId := middleware.GetUserIDFromContext(ctx)

	establishment, err := c.establishmentService.GetEstablishmentByAdminID(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}

	overdueAccounts, err := c.creditAccountService.GetOverdueCreditAccounts(establishment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, overdueAccounts)
}

// ProcessPurchase godoc
// @Summary      Process Purchase
// @Description  Processes a purchase on a client's credit account.
// @Tags         Credit Accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        creditAccountID path int true "Credit Account ID"
// @Param        purchase        body      request.CreateTransactionRequest  true  "Purchase details"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/{creditAccountID}/purchases [post]
func (c *CreditAccountController) ProcessPurchase(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	var req request.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure transaction type is purchase
	if req.TransactionType != enums.Purchase {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid transaction type for purchase"})
		return
	}

	// Only admins can process purchases
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can process purchases"})
		return
	}

	// Additional validation if needed...

	err = c.creditAccountService.ProcessPurchase(uint(creditAccountID), req.Amount, req.Description)
	if err != nil {
		// Handle different error types appropriately (e.g., validation errors, insufficient credit, etc.)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Purchase processed successfully"})
}

// ProcessPayment godoc
// @Summary      Process Payment
// @Description  Processes a payment towards a client's credit account.
// @Tags         Credit Accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        creditAccountID path int true "Credit Account ID"
// @Param        payment        body      request.CreateTransactionRequest  true  "Payment details"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/{creditAccountID}/payments [post]
func (c *CreditAccountController) ProcessPayment(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	var req request.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Ensure transaction type is payment
	if req.TransactionType != enums.Payment {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid transaction type for payment"})
		return
	}

	// Only admins can process payments (for now)
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can process payments"})
		return
	}

	// Additional validation if needed...

	err = c.creditAccountService.ProcessPayment(uint(creditAccountID), req.Amount, req.Description)
	if err != nil {
		// Handle different error types appropriately (e.g., validation errors, insufficient funds, etc.)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Payment processed successfully"})
}

// GetAdminDebtSummary godoc
// @Summary      Get Admin Debt Summary
// @Description  Retrieves a summary of all client debts for an establishment. Only Admins can access this endpoint.
// @Tags         Credit Accounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {array}  response.AdminDebtSummary
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/debt-summary [get]
func (c *CreditAccountController) GetAdminDebtSummary(ctx *gin.Context) {
	// Ensure the authenticated user is an ADMIN
	if middleware.GetUserRoleFromContext(ctx) != enums.ADMIN {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only admins can access this endpoint"})
		return
	}

	userId := middleware.GetUserIDFromContext(ctx)

	establishment, err := c.establishmentService.GetEstablishmentByAdminID(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}
	fmt.Println(establishment.ID) // Debugging line

	summary, err := c.creditAccountService.GetAdminDebtSummary(establishment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// UpdateCreditAccountByClientID godoc
// @Summary      Update Credit Account by Client ID
// @Description  Updates an existing credit account by client ID. Only Admins can update credit accounts.
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
func (c *CreditAccountController) UpdateCreditAccountByClientID(ctx *gin.Context) {
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account not found for this client"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccountResponse)
}
