package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreditAccountController handles API requests related to credit accounts.
type CreditAccountController struct {
	creditAccountService service.CreditAccountService
}

// NewCreditAccountController creates a new CreditAccountController.
func NewCreditAccountController(creditAccountService service.CreditAccountService) *CreditAccountController {
	return &CreditAccountController{creditAccountService: creditAccountService}
}

// CreateCreditAccount godoc
// @Summary Create a new credit account
// @Description Creates a new credit account for a client.
// @Tags CreditAccounts
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccount body request.CreateCreditAccountRequest true "Credit account details"
// @Success 201 {object} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts [post]
func (c *CreditAccountController) CreateCreditAccount(ctx *gin.Context) {
	var req request.CreateCreditAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	creditAccount, err := c.creditAccountService.CreateCreditAccount(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, creditAccount)
}

// GetCreditAccountByID godoc
// @Summary Get credit account by ID
// @Description Retrieves a credit account by its ID.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Credit Account ID"
// @Success 200 {object} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{id} [get]
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

// UpdateCreditAccount godoc
// @Summary Update a credit account
// @Description Updates an existing credit account by its ID.
// @Tags CreditAccounts
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Credit Account ID"
// @Param creditAccount body request.UpdateCreditAccountRequest true "Updated credit account details"
// @Success 200 {object} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{id} [put]
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
// @Summary Delete a credit account
// @Description Deletes a credit account by its ID.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Credit Account ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{id} [delete]
func (c *CreditAccountController) DeleteCreditAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
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
// @Summary Get credit accounts by establishment ID
// @Description Retrieves all credit accounts associated with an establishment.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param establishmentID path int true "Establishment ID"
// @Success 200 {array} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishmentID}/credit-accounts [get]
func (c *CreditAccountController) GetCreditAccountsByEstablishmentID(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	creditAccounts, err := c.creditAccountService.GetCreditAccountsByEstablishmentID(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccounts)
}

// GetCreditAccountsByClientID godoc
// @Summary Get credit accounts by client ID
// @Description Retrieves all credit accounts associated with a client.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param clientID path int true "Client ID"
// @Success 200 {array} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /clients/{clientID}/credit-accounts [get]
func (c *CreditAccountController) GetCreditAccountsByClientID(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("clientID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	creditAccounts, err := c.creditAccountService.GetCreditAccountsByClientID(uint(clientID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccounts)
}

// ApplyInterestToAllAccounts godoc
// @Summary Apply interest to all accounts
// @Description Applies interest to all eligible credit accounts within an establishment.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param establishmentID path int true "Establishment ID"
// @Success 200 "Interest applied successfully"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishmentID}/credit-accounts/apply-interest [post]
func (c *CreditAccountController) ApplyInterestToAllAccounts(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	if err := c.creditAccountService.ApplyInterestToAllAccounts(uint(establishmentID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interest applied successfully"})
}

// ApplyLateFeesToAllAccounts godoc
// @Summary Apply late fees to all accounts
// @Description Applies late fees to all eligible credit accounts within an establishment.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param establishmentID path int true "Establishment ID"
// @Success 200 "Late fees applied successfully"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishmentID}/credit-accounts/apply-late-fees [post]
func (c *CreditAccountController) ApplyLateFeesToAllAccounts(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	if err := c.creditAccountService.ApplyLateFeesToAllAccounts(uint(establishmentID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Late fees applied successfully"})
}

// GetAdminDebtSummary godoc
// @Summary Get admin debt summary
// @Description Retrieves a summary of debts owed to an establishment.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param establishmentID path int true "Establishment ID"
// @Success 200 {array} response.AdminDebtSummary
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishmentID}/credit-accounts/debt-summary [get]
func (c *CreditAccountController) GetAdminDebtSummary(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	summary, err := c.creditAccountService.GetAdminDebtSummary(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, summary)
}

// ProcessPurchase godoc
// @Summary Process a purchase
// @Description Processes a purchase transaction on a credit account.
// @Tags CreditAccounts
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Param purchase body request.CreateTransactionRequest true "Purchase details"
// @Success 201 "Purchase processed successfully"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/purchases [post]
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

	// You may want to add additional validation here,
	// such as checking if the transaction type is Purchase.

	if err := c.creditAccountService.ProcessPurchase(uint(creditAccountID), req.Amount, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Purchase processed successfully"})
}

// ProcessPayment godoc
// @Summary Process a payment
// @Description Processes a payment transaction on a credit account.
// @Tags CreditAccounts
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Param payment body request.CreateTransactionRequest true "Payment details"
// @Success 201 "Payment processed successfully"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/payments [post]
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

	

	if err := c.creditAccountService.ProcessPayment(uint(creditAccountID), req.Amount, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Payment processed successfully"})
}

// CreateCreditRequest godoc
// @Summary Create a credit request
// @Description Creates a new credit request for a client.
// @Tags CreditRequests
// @Accept json
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditRequest body request.CreateCreditRequest true "Credit request details"
// @Success 201 {object} response.CreditRequestResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-requests [post]
func (c *CreditAccountController) CreateCreditRequest(ctx *gin.Context) {
	var req request.CreateCreditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	creditRequest, err := c.creditAccountService.CreateCreditRequest(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, creditRequest)
}

// GetCreditRequestByID godoc
// @Summary Get credit request by ID
// @Description Retrieves a credit request by its ID.
// @Tags CreditRequests
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Credit Request ID"
// @Success 200 {object} response.CreditRequestResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-requests/{id} [get]
func (c *CreditAccountController) GetCreditRequestByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit request ID"})
		return
	}

	creditRequest, err := c.creditAccountService.GetCreditRequestByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit request not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditRequest)
}

// GetCreditRequestsByClientID godoc
// @Summary Get credit requests by Client ID
// @Description Retrieves all credit requests belonging to a specific client.
// @Tags CreditRequests
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param client_id path int true "Client ID"
// @Success 200 {array} response.CreditRequestResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse 
// @Failure 500 {object} response.ErrorResponse
// @Router /clients/{client_id}/credit-requests [get]
func (c *CreditAccountController) GetCreditRequestsByClientID(ctx *gin.Context) {
    clientID, err := strconv.Atoi(ctx.Param("client_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
        return
    }

    creditRequests, err := c.creditAccountService.GetCreditRequestsByClientID(uint(clientID))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "No credit requests found for this client"}) 
            return
        }
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, creditRequests)
}

// GetCreditRequestsByEstablishmentID godoc
// @Summary Get credit requests by Establishment ID
// @Description Retrieves all credit requests associated with an establishment.
// @Tags CreditRequests
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param establishment_id path int true "Establishment ID"
// @Success 200 {array} response.CreditRequestResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishment_id}/credit-requests [get]
func (c *CreditAccountController) GetCreditRequestsByEstablishmentID(ctx *gin.Context) {
    establishmentID, err := strconv.Atoi(ctx.Param("establishment_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
        return
    }

    creditRequests, err := c.creditAccountService.GetCreditRequestsByEstablishmentID(uint(establishmentID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, creditRequests)
}


// ApproveCreditRequest godoc
// @Summary Approve credit request
// @Description Approves a pending credit request and creates a credit account.
// @Tags CreditRequests
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Credit Request ID"
// @Success 200 {object} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-requests/{id}/approve [put]
func (c *CreditAccountController) ApproveCreditRequest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit request ID"})
		return
	}

	// Assuming you have middleware to retrieve the authenticated admin\'s ID
	adminID := middleware.GetUserIDFromContext(ctx)

	creditAccount, err := c.creditAccountService.ApproveCreditRequest(uint(id), adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit request not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditAccount)
}

// RejectCreditRequest godoc
// @Summary Reject credit request
// @Description Rejects a pending credit request.
// @Tags CreditRequests
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param id path int true "Credit Request ID"
// @Success 200 "Credit request rejected successfully"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-requests/{id}/reject [put]
func (c *CreditAccountController) RejectCreditRequest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit request ID"})
		return
	}

	// Assuming you have middleware to retrieve the authenticated admin's ID
	adminID := middleware.GetUserIDFromContext(ctx)

	if err := c.creditAccountService.RejectCreditRequest(uint(id), adminID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit request not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Credit request rejected successfully"})
}

// GetPendingCreditRequests godoc
// @Summary Get pending credit requests
// @Description Retrieves all pending credit requests for an establishment.
// @Tags CreditRequests
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param establishmentID path int true "Establishment ID"
// @Success 200 {array} response.CreditRequestResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /establishments/{establishmentID}/credit-requests/pending [get]
func (c *CreditAccountController) GetPendingCreditRequests(ctx *gin.Context) {
	establishmentID, err := strconv.Atoi(ctx.Param("establishmentID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid establishment ID"})
		return
	}

	creditRequests, err := c.creditAccountService.GetPendingCreditRequests(uint(establishmentID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, creditRequests)
}

// UpdateCreditRequestDueDate godoc
// @Summary Update credit request due date
// @Description Updates the due date of a credit request.
// @Tags CreditRequests
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Credit Request ID"
// @Param dueDate body request.UpdateCreditRequestDueDate true "New due date"
// @Success 200 {object} response.CreditRequestResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-requests/{id}/due-date [put] 
func (c *CreditAccountController) UpdateCreditRequestDueDate(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit request ID"})
        return
    }

    var req request.UpdateCreditRequestDueDate
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
        return
    }

    creditRequest, err := c.creditAccountService.UpdateCreditRequestDueDate(uint(id), req.DueDate)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit request not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, creditRequest)
}


// AssignCreditAccountToClient godoc
// @Summary Assign credit account to client
// @Description Assigns an existing credit account to a client.
// @Tags CreditAccounts
// @Produce json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param creditAccountID path int true "Credit Account ID"
// @Param clientID path int true "Client ID"
// @Success 200 {object} response.CreditAccountResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /credit-accounts/{creditAccountID}/clients/{clientID} [put]
func (c *CreditAccountController) AssignCreditAccountToClient(ctx *gin.Context) {
	creditAccountID, err := strconv.Atoi(ctx.Param("creditAccountID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid credit account ID"})
		return
	}

	clientID, err := strconv.Atoi(ctx.Param("clientID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	updatedCreditAccount, err := c.creditAccountService.AssignCreditAccountToClient(uint(creditAccountID), uint(clientID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Credit account or client not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedCreditAccount)
}


// GetClientAccountStatement godoc
// @Summary      Get Client Account Statement
// @Description  Retrieves a client's account statement within a specified date range.
// @Tags         CreditAccounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        clientID       path        int     true  "Client ID"
// @Param        startDate      query       string  false "Start date (YYYY-MM-DD)"
// @Param        endDate        query       string  false "End date (YYYY-MM-DD)"
// @Success      200  {object}  response.AccountStatementResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/clients/{clientID}/statement [get]
func (c *CreditAccountController) GetClientAccountStatement(ctx *gin.Context) {
    clientID, err := strconv.Atoi(ctx.Param("clientID"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
        return
    }

    startDateStr := ctx.Query("startDate") 
    endDateStr := ctx.Query("endDate")   

    var startDate, endDate time.Time
    if startDateStr != "" {
        startDate, err = time.Parse("2006-01-02", startDateStr) // Adjust date format as needed
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

    statement, err := c.creditAccountService.GetClientAccountStatement(uint(clientID), startDate, endDate)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, statement)
}

// GetClientAccountHistory godoc
// @Summary      Get Client Account History
// @Description  Retrieves the complete transaction history for a client's credit account.
// @Tags         CreditAccounts
// @Produce      json
// @Param        Authorization  header      string  true  "Bearer {token}"
// @Param        clientID       path        int     true  "Client ID"
// @Success      200  {object}  response.AccountStatementResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /credit-accounts/clients/{clientID}/history [get]
func (c *CreditAccountController) GetClientAccountHistory(ctx *gin.Context) {
    clientID, err := strconv.Atoi(ctx.Param("clientID"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid client ID"})
        return
    }

    history, err := c.creditAccountService.GetClientAccountHistory(uint(clientID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, history)
}
