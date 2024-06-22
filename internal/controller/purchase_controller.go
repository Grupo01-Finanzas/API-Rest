package controller

import (
	"net/http"
	"time"

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
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /purchases [post]
func (c *PurchaseController) CreatePurchase(ctx *gin.Context) {
	var req request.CreatePurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID and role from context
	userID := middleware.GetUserIDFromContext(ctx)
	userRole := middleware.GetUserRoleFromContext(ctx)

	// Only Clients can make purchases
	if userRole != enums.CLIENT {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Only clients can make purchases"})
		return
	}

	// Validate credit type
	if req.CreditType != enums.ShortTerm && req.CreditType != enums.LongTerm {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit type"})
		return
	}

	err := c.purchaseService.ProcessPurchase(userID, req.EstablishmentID, req.ProductIDs, req.CreditType, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Purchase created successfully"})
}

// GetClientBalance godoc
// @Summary      Get Client Balance
// @Description  Gets the current balance of the authenticated client's credit account.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  response.ClientBalanceResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/balance [get]
func (c *PurchaseController) GetClientBalance(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	balance, err := c.purchaseService.GetClientBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp := response.ClientBalanceResponse{
		ClientID:       userID,
		CurrentBalance: balance,
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetClientTransactions godoc
// @Summary      Get Client Transactions
// @Description  Gets the transaction history of the authenticated client.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {array}   response.TransactionResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/transactions [get]
func (c *PurchaseController) GetClientTransactions(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	transactions, err := c.purchaseService.GetClientTransactions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// GetClientOverdueBalance godoc
// @Summary      Get Client Overdue Balance
// @Description  Gets the overdue balance of the authenticated client's credit account.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  map[string]float64
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/overdue-balance [get]
func (c *PurchaseController) GetClientOverdueBalance(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	overdueBalance, err := c.purchaseService.GetClientOverdueBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"overdue_balance": overdueBalance})
}

// GetClientInstallments godoc
// @Summary      Get Client Installments
// @Description  Gets the installments of the authenticated client's credit account.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {array}   response.InstallmentResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/installments [get]
func (c *PurchaseController) GetClientInstallments(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	installments, err := c.purchaseService.GetClientInstallments(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installments)
}

// GetClientCreditAccount godoc
// @Summary      Get Client Credit Account
// @Description  Gets the credit account details of the authenticated client.
// @Tags         Clients
// @Accept       json
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  response.CreditAccountResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/credit-account [get]
func (c *PurchaseController) GetClientCreditAccount(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	creditAccount, err := c.purchaseService.GetClientCreditAccount(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccount)
}

// GetClientAccountSummary godoc
// @Summary      Get Client Account Summary
// @Description  Retrieves a summary of the client's account, including transactions, payments, debts, and interest.
// @Tags         Clients
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Success      200  {object}  response.AccountSummaryResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/account-summary [get]
func (c *PurchaseController) GetClientAccountSummary(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)

	summary, err := c.purchaseService.GetClientAccountSummary(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// GetClientAccountStatement godoc
// @Summary      Get Client Account Statement
// @Description  Retrieves an account statement for the client within a specified date range.
// @Tags         Clients
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        startDate      query       string  false "Start date (YYYY-MM-DD)"
// @Param        endDate        query       string  false "End date (YYYY-MM-DD)"
// @Success      200  {object}  response.AccountStatementResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/account-statement [get]
func (c *PurchaseController) GetClientAccountStatement(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid start date format"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid end date format"})
			return
		}
	}

	statement, err := c.purchaseService.GetClientAccountStatement(userID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, statement)
}

// GetClientAccountStatementPDF godoc
// @Summary      Get Client Account Statement (PDF)
// @Description  Generates and downloads a PDF account statement for the client within a specified date range.
// @Tags         Clients
// @Produce      application/pdf
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        startDate      query       string  false "Start date (YYYY-MM-DD)"
// @Param        endDate        query       string  false "End date (YYYY-MM-DD)"
// @Success      200  {file}   application/pdf  "PDF Account Statement"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /clients/me/account-statement/pdf [get]
func (c *PurchaseController) GetClientAccountStatementPDF(ctx *gin.Context) {
	userID := middleware.GetUserIDFromContext(ctx)
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid start date format"})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid end date format"})
			return
		}
	}

	// Get the PDF data from the service
	pdfBytes, err := c.purchaseService.GenerateClientAccountStatementPDF(userID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error generating PDF: " + err.Error()})
		return
	}

	// Set headers for PDF download
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=account_statement.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
