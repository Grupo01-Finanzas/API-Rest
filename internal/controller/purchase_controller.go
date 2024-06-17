package controller

import (
	"net/http"
	"strconv"

	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/service"
	"github.com/gin-gonic/gin"
)

// PurchaseController handles endpoints related to purchases
type PurchaseController struct {
	purchaseService service.PurchaseService
}

// NewPurchaseController creates a new instance of PurchaseController
func NewPurchaseController(purchaseService service.PurchaseService) *PurchaseController {
	return &PurchaseController{purchaseService: purchaseService}
}

// CreatePurchase godoc
// @Summary      Create a Purchase
// @Description  Processes a product purchase by a user.
// @Tags         Purchases
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        purchase         body      request.CreatePurchaseRequest  true  "Purchase Data"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /purchases [post]
func (c *PurchaseController) CreatePurchase(ctx *gin.Context) {
	var req request.CreatePurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserIDFromContext(ctx)

	// Validate credit type
	if req.CreditType != enums.ShortTerm && req.CreditType != enums.LongTerm {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credit type"})
		return
	}

	err := c.purchaseService.ProcessPurchase(userID, req.EstablishmentID, req.ProductIDs, req.CreditType, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Purchase created successfully"})
}

// GetClientBalance godoc
// @Summary      Get Client Balance
// @Description  Gets the current balance of a client's credit account.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  response.ClientBalanceResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clients/{id}/balance [get]
func (c *PurchaseController) GetClientBalance(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID"})
		return
	}

	// Get client balance from the service
	balance, err := c.purchaseService.GetClientBalance(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response.ClientBalanceResponse{
		ClientID:       uint(clientID),
		CurrentBalance: balance,
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetClientTransactions godoc
// @Summary      Get Client Transactions
// @Description  Gets the transaction history of a client.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {array}   response.TransactionResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clients/{id}/transactions [get]
func (c *PurchaseController) GetClientTransactions(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID"})
		return
	}

	// Get client transactions from the service
	transactions, err := c.purchaseService.GetClientTransactions(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// GetClientOverdueBalance godoc
// @Summary      Get Client Overdue Balance
// @Description  Gets the overdue balance of a client's credit account.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  map[string]float64
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clients/{id}/overdue-balance [get]
func (c *PurchaseController) GetClientOverdueBalance(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID"})
		return
	}

	// Get client overdue balance from the service
	overdueBalance, err := c.purchaseService.GetClientOverdueBalance(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"overdue_balance": overdueBalance})
}

// GetClientInstallments godoc
// @Summary      Get Client Installments
// @Description  Gets the installments of a client's credit account.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {array}   response.InstallmentResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clients/{id}/installments [get]
func (c *PurchaseController) GetClientInstallments(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID"})
		return
	}

	// Get client installments from the service
	installments, err := c.purchaseService.GetClientInstallments(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installments)
}

// GetClientCreditAccount godoc
// @Summary      Get Client Credit Account
// @Description  Gets the credit account details of a client.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  response.CreditAccountResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /clients/{id}/credit-account [get]
func (c *PurchaseController) GetClientCreditAccount(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID"})
		return
	}

	// Get client credit account from the service
	creditAccount, err := c.purchaseService.GetClientCreditAccount(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccount)
}
