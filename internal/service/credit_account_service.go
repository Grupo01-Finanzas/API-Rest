package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"gorm.io/gorm"
)

// CreditAccountService defines methods for managing credit accounts.
type CreditAccountService interface {
	CreateCreditAccount(req request.CreateCreditAccountRequest) (*response.CreditAccountResponse, error)
	GetCreditAccountByID(id uint) (*response.CreditAccountResponse, error)
	UpdateCreditAccount(id uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error)
	DeleteCreditAccount(id uint) error
	GetCreditAccountsByEstablishmentID(establishmentID uint) ([]response.CreditAccountResponse, error)
	GetCreditAccountsByClientID(clientID uint) ([]response.CreditAccountResponse, error)
	ApplyInterestToAllAccounts(establishmentID uint) error
	ApplyLateFeesToAllAccounts(establishmentID uint) error
	GetAdminDebtSummary(establishmentID uint) ([]response.AdminDebtSummary, error)
	ProcessPurchase(creditAccountID uint, amount float64, description string) error
	ProcessPayment(creditAccountID uint, amount float64, description string) error

	CreateCreditRequest(req request.CreateCreditRequest) (*response.CreditRequestResponse, error)
	GetCreditRequestByID(id uint) (*response.CreditRequestResponse, error)
	ApproveCreditRequest(creditRequestID uint, adminID uint) (*response.CreditAccountResponse, error)
	RejectCreditRequest(creditRequestID uint, adminID uint) error
	GetPendingCreditRequests(establishmentID uint) ([]response.CreditRequestResponse, error)
	AssignCreditAccountToClient(creditAccountID, clientID uint) (*response.CreditAccountResponse, error)
	GetNumberOfDues(account entities.CreditAccount) int
	CalculateInterest(creditAccount entities.CreditAccount) float64

	GetClientAccountStatement(clientID uint, startDate, endDate time.Time) ([]*response.AccountStatementResponse, error)
	CalculateDueDate(account entities.CreditAccount) (time.Time, error)
	GetClientAccountHistory(clientID uint) (*response.AccountStatementResponse, error)
}

type creditAccountService struct {
	creditAccountRepo repository.CreditAccountRepository
	transactionRepo   repository.TransactionRepository
	clientRepo        repository.ClientRepository
	establishmentRepo repository.EstablishmentRepository
	installmentRepo   repository.InstallmentRepository
}

// NewCreditAccountService creates a new instance of CreditAccountService.
func NewCreditAccountService(creditAccountRepo repository.CreditAccountRepository, transactionRepo repository.TransactionRepository, clientRepo repository.ClientRepository, establishmentRepo repository.EstablishmentRepository, installmentRepo repository.InstallmentRepository) CreditAccountService {
	return &creditAccountService{
		creditAccountRepo: creditAccountRepo,
		transactionRepo:   transactionRepo,
		clientRepo:        clientRepo,
		establishmentRepo: establishmentRepo,
		installmentRepo:   installmentRepo,
	}
}

func (s *creditAccountService) CreateCreditAccount(req request.CreateCreditAccountRequest) (*response.CreditAccountResponse, error) {
	// 1. Check if the client exists
	clientResponse, err := s.clientRepo.GetClientByID(req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client: %w", err)
	}
	if clientResponse == nil {
		return nil, fmt.Errorf("client with ID %d not found", req.ClientID)
	}

	// 2. Check if the establishment exists
	establishmentResponse, err := s.establishmentRepo.GetByID(req.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving establishment: %w", err)
	}
	if establishmentResponse == nil {
		return nil, fmt.Errorf("establishment with ID %d not found", req.EstablishmentID)
	}

	// 3. Check if a credit account already exists for this client and establishment
	exists, err := s.creditAccountRepo.ExistsByClientAndEstablishment(req.ClientID, req.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing credit account: %w", err)
	}
	if exists {
		return nil, errors.New("a credit account already exists for this client and establishment")
	}

	// 4. Create the credit account (assuming the establishment has permission)
	creditAccount := entities.CreditAccount{
		EstablishmentID:         req.EstablishmentID,
		ClientID:                req.ClientID,
		CreditLimit:             req.CreditLimit,
		MonthlyDueDate:          req.MonthlyDueDate,
		InterestRate:            req.InterestRate,
		InterestType:            req.InterestType,
		CreditType:              req.CreditType,
		GracePeriod:             req.GracePeriod,
		IsBlocked:               false, // Initially not blocked
		LastInterestAccrualDate: time.Now(),
		CurrentBalance:          0.0, // Initial balance is zero
		LateFeeRuleID:           req.LateFeeRuleID,
	}

	// 5. Create the credit account using the repository
	if err := s.creditAccountRepo.Create(&creditAccount); err != nil {
		return nil, fmt.Errorf("error creating credit account: %w", err)
	}

	return creditAccountToResponse(&creditAccount), nil
}

func (s *creditAccountService) GetCreditAccountByID(id uint) (*response.CreditAccountResponse, error) {
	return s.creditAccountRepo.GetByID(id)
}

func (s *creditAccountService) UpdateCreditAccount(id uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error) {
	// You can add additional business logic here before updating the account
	// For example, validation, permissions, etc.
	return s.creditAccountRepo.Update(id, req)
}

func (s *creditAccountService) DeleteCreditAccount(id uint) error {
	return s.creditAccountRepo.Delete(id)
}

func (s *creditAccountService) GetCreditAccountsByEstablishmentID(establishmentID uint) ([]response.CreditAccountResponse, error) {
	creditAccounts, err := s.creditAccountRepo.GetByEstablishmentID(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit accounts: %w", err)
	}

	creditAccountResponses := make([]response.CreditAccountResponse, len(creditAccounts))
	copy(creditAccountResponses, creditAccounts)
	return creditAccountResponses, nil
}

func (s *creditAccountService) GetCreditAccountsByClientID(clientID uint) ([]response.CreditAccountResponse, error) {
	return s.creditAccountRepo.GetByClientID(clientID)
}

// ApplyInterestToAllAccounts applies interest to all eligible credit accounts within an establishment.
func (s *creditAccountService) ApplyInterestToAllAccounts(establishmentID uint) error {
	creditAccounts, err := s.creditAccountRepo.GetByEstablishmentID(establishmentID)
	if err != nil {
		return fmt.Errorf("error retrieving credit accounts: %w", err)
	}

	for _, account := range creditAccounts {
		if err := s.creditAccountRepo.ApplyInterest(account.ID); err != nil {
			return fmt.Errorf("error applying interest to account %d: %w", account.ID, err)
		}
	}

	return nil
}

// ApplyLateFeesToAllAccounts applies late fees to all eligible credit accounts within an establishment.
func (s *creditAccountService) ApplyLateFeesToAllAccounts(establishmentID uint) error {
	overdueAccounts, err := s.creditAccountRepo.GetOverdueAccounts(establishmentID)
	if err != nil {
		return fmt.Errorf("error retrieving overdue accounts: %w", err)
	}

	for _, account := range overdueAccounts {
		if err := s.creditAccountRepo.ApplyLateFee(account.ID); err != nil {
			return fmt.Errorf("error applying late fee to account %d: %w", account.ID, err)
		}
	}

	return nil
}

// GetAdminDebtSummary retrieves a summary of debts owed to an establishment.
func (s *creditAccountService) GetAdminDebtSummary(establishmentID uint) ([]response.AdminDebtSummary, error) {
    creditAccounts, err := s.creditAccountRepo.GetByEstablishmentID(establishmentID)
    if err != nil {
        return nil, fmt.Errorf("error retrieving credit accounts: %w", err)
    }

    summary := make([]response.AdminDebtSummary, 0, len(creditAccounts))
    for _, account := range creditAccounts {
        client, err := s.clientRepo.GetClientByID(account.ClientID)
        if err != nil {
            return nil, fmt.Errorf("error retrieving client: %w", err)
        }

        dueDate, err := s.CalculateDueDate(*responseToCreditAccount(&account)) // Get both return values
        if err != nil {
            return nil, fmt.Errorf("error calculating due date: %w", err)  // Handle the error
        }

        summaryItem := response.AdminDebtSummary{
            ClientID:       account.ClientID,
            ClientName:     client.User.Name,
            CreditType:     string(account.CreditType),
            InterestRate:   account.InterestRate,
            NumberOfDues:   s.GetNumberOfDues(*responseToCreditAccount(&account)), 
            CurrentBalance: account.CurrentBalance,
            DueDate:        dueDate, 
        }

        summary = append(summary, summaryItem)
    }

    return summary, nil 
}

// ProcessPurchase processes a purchase on a credit account.
func (s *creditAccountService) ProcessPurchase(creditAccountID uint, amount float64, description string) error {
    // Check for overdue balance
    overdueBalance, err := s.creditAccountRepo.GetOverdueBalance(creditAccountID)
    if err != nil {
        return fmt.Errorf("error checking overdue balance: %w", err)
    }

    if overdueBalance > 0 {
        return fmt.Errorf("cannot process purchase: client has an overdue balance of %.2f", overdueBalance) // Include balance in error message
    }

    // Proceed with purchase if no overdue balance
    return s.creditAccountRepo.ProcessPurchase(creditAccountID, amount, description) 
}

// ProcessPayment processes a payment towards a credit account.
func (s *creditAccountService) ProcessPayment(creditAccountID uint, amount float64, description string) error {
	// The service method now simply calls the repository method
	return s.creditAccountRepo.ProcessPayment(creditAccountID, amount, description)
}

// Helper functions for calculations

func (s *creditAccountService) GetNumberOfDues(account entities.CreditAccount) int {
	if account.CreditType != enums.LongTerm {
		return 0 // Number of dues is not applicable for ShortTerm credit
	}

	// 1. Retrieve installments from the database
	installments, err := s.installmentRepo.GetByCreditAccountID(account.ID)
	if err != nil {
		// Handle error appropriately (e.g., log the error and return 0)
		return 0
	}

	// 2. Count the number of installments
	numberOfDues := len(installments)

	return numberOfDues
}

func (s *creditAccountService) CalculateInterest(creditAccount entities.CreditAccount) float64 {
	var interest float64
	today := time.Now()

	if creditAccount.CreditType == enums.ShortTerm {
		// Short-Term Interest Calculation

		// Calculate the number of days since the last interest accrual
		daysSinceLastAccrual := today.Sub(creditAccount.LastInterestAccrualDate).Hours() / 24

		// Check if it's time to accrue interest (at least a month has passed)
		if daysSinceLastAccrual >= daysInMonth(creditAccount.LastInterestAccrualDate.Month(), creditAccount.LastInterestAccrualDate.Year()) {
			if creditAccount.InterestType == enums.Nominal {
				// Nominal Interest: Interest = Principal * (Rate/100) * (Time in years)
				interest = creditAccount.CurrentBalance * (creditAccount.InterestRate / 100) * (daysSinceLastAccrual / 365.0)
			} else if creditAccount.InterestType == enums.Effective {
				// Effective Interest: Interest = Principal * ((1 + Rate/100)^(Time in years) - 1)
				interest = creditAccount.CurrentBalance * (math.Pow(1+(creditAccount.InterestRate/100), daysSinceLastAccrual/365.0) - 1)
			}
		}

	} else { // LongTerm
		// Long-Term (Installment) Interest Calculation
		installments, err := s.installmentRepo.GetByCreditAccountID(creditAccount.ID)
		if err != nil {
			
			return 0
		}

		for _, installment := range installments {
			// Only calculate interest on pending installments
			if installment.Status == enums.Pending {
				// Calculate days until the installment's due date
				daysUntilDueDate := time.Until(installment.DueDate).Hours() / 24

				// Check if the due date is in the future
				if daysUntilDueDate > 0 {
					if creditAccount.InterestType == enums.Nominal {
						// Nominal Interest for Installment
						interest += installment.Amount * (creditAccount.InterestRate / 100) * (daysUntilDueDate / 365.0)
					} else if creditAccount.InterestType == enums.Effective {
						// Effective Interest for Installment
						interest += installment.Amount * (math.Pow(1+(creditAccount.InterestRate/100), daysUntilDueDate/365.0) - 1)
					}
				}
			}
		}
	}

	return interest
}

// Helper function to get the number of days in a month
func daysInMonth(month time.Month, year int) float64 {
	return float64(time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day())
}

func (s *creditAccountService) CreateCreditRequest(req request.CreateCreditRequest) (*response.CreditRequestResponse, error) {
	// 1. Check if a credit account already exists for the client and establishment
	exists, err := s.creditAccountRepo.ExistsByClientAndEstablishment(req.ClientID, req.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing credit account: %w", err)
	}
	if exists {
		return nil, errors.New("a credit account already exists for this client and establishment")
	}

	// 2. Create the credit request
	creditRequest := entities.CreditRequest{
		ClientID:             req.ClientID,
		EstablishmentID:      req.EstablishmentID,
		RequestedCreditLimit: req.RequestedCreditLimit,
		MonthlyDueDate:       req.MonthlyDueDate,
		InterestType:         req.InterestType,
		CreditType:           req.CreditType,
		GracePeriod:          req.GracePeriod,
		Status:               entities.Pending,
	}

	if err := s.creditAccountRepo.CreateCreditRequest(&creditRequest); err != nil { // Use the repository method
		return nil, fmt.Errorf("error creating credit request: %w", err)
	}

	// 3. Map creditRequest to CreditRequestResponse and return
	return getCreditRequestResponse(&creditRequest), nil
}

func (s *creditAccountService) GetCreditRequestByID(id uint) (*response.CreditRequestResponse, error) {
	creditRequest, err := s.creditAccountRepo.GetCreditRequestByID(id) // Use the repository method
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit request: %w", err)
	}

	return getCreditRequestResponse(creditRequest), nil
}

func (s *creditAccountService) ApproveCreditRequest(creditRequestID uint, adminID uint) (*response.CreditAccountResponse, error) {
	// 1. Retrieve the credit request using the repository (now preloads data)
	creditRequest, err := s.creditAccountRepo.GetCreditRequestByID(creditRequestID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit request: %w", err)
	}

	// 2. Check if the admin belongs to the establishment
	if creditRequest.Establishment.Admin.UserID != adminID {
		return nil, errors.New("admin does not have permission to approve this request")
	}

	// 3. Check if the request is already approved or rejected
	if creditRequest.Status != entities.Pending {
		return nil, fmt.Errorf("credit request is already %s", creditRequest.Status)
	}

	// 4. Approve the credit request in the repository (no more preloading needed here)
	creditAccountResponse, err := s.creditAccountRepo.ApproveCreditRequest(creditRequest)
	if err != nil {
		return nil, fmt.Errorf("error approving credit request: %w", err)
	}

	return creditAccountResponse, nil
}

func (s *creditAccountService) RejectCreditRequest(creditRequestID uint, adminID uint) error {
	// 1. Retrieve the credit request
	creditRequest, err := s.creditAccountRepo.GetCreditRequestByID(creditRequestID)
	if err != nil {
		return fmt.Errorf("error retrieving credit request: %w", err)
	}

	// 2. Check if the admin belongs to the establishment
	if creditRequest.Establishment.Admin.UserID != adminID {
		return errors.New("admin does not have permission to reject this request")
	}

	// 3. Check if the request is already approved or rejected
	if creditRequest.Status != entities.Pending {
		return fmt.Errorf("credit request is already %s", creditRequest.Status)
	}

	// 4. Update the credit request status
	now := time.Now()
	creditRequest.Status = entities.Rejected
	creditRequest.RejectedAt = &now

	// 5. Use the repository method to update the request
	if err := s.creditAccountRepo.UpdateCreditRequest(creditRequest); err != nil {
		return fmt.Errorf("error updating credit request status: %w", err)
	}

	return nil
}

func (s *creditAccountService) GetPendingCreditRequests(establishmentID uint) ([]response.CreditRequestResponse, error) {
	// 1. Retrieve credit requests using the repository method
	creditRequests, err := s.creditAccountRepo.GetPendingCreditRequests(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving pending credit requests: %w", err)
	}

	// 2. Convert entities to responses
	var creditRequestResponses []response.CreditRequestResponse
	for _, creditRequest := range creditRequests {
		creditRequestResponses = append(creditRequestResponses, *getCreditRequestResponse(&creditRequest))
	}

	return creditRequestResponses, nil
}

func getCreditRequestResponse(creditRequest *entities.CreditRequest) *response.CreditRequestResponse {
	return &response.CreditRequestResponse{
		ID:                   creditRequest.ID,
		ClientID:             creditRequest.ClientID,
		EstablishmentID:      creditRequest.EstablishmentID,
		RequestedCreditLimit: creditRequest.RequestedCreditLimit,
		MonthlyDueDate:       creditRequest.MonthlyDueDate,
		InterestType:         creditRequest.InterestType,
		CreditType:           creditRequest.CreditType,
		GracePeriod:          creditRequest.GracePeriod,
		Status:               creditRequest.Status,
		ApprovedAt:           creditRequest.ApprovedAt,
		RejectedAt:           creditRequest.RejectedAt,
		CreatedAt:            creditRequest.CreatedAt,
		UpdatedAt:            creditRequest.UpdatedAt,
	}
}

// AssignCreditAccountToClient assigns an existing credit account to a client.
func (s *creditAccountService) AssignCreditAccountToClient(creditAccountID, clientID uint) (*response.CreditAccountResponse, error) {
	// Call the repository method to perform the assignment
	err := s.creditAccountRepo.AssignCreditAccountToClient(creditAccountID, clientID)
	if err != nil {
		return nil, fmt.Errorf("error assigning credit account to client: %w", err)
	}

	// If successful, retrieve the updated credit account and return the response
	updatedCreditAccount, err := s.creditAccountRepo.GetByID(creditAccountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated credit account: %w", err)
	}

	return updatedCreditAccount, nil
}

func creditAccountToResponse(creditAccount *entities.CreditAccount) *response.CreditAccountResponse {
	clientResponse := NewClientResponse(creditAccount.Client)
	return &response.CreditAccountResponse{
		ID:                      creditAccount.ID,
		EstablishmentID:         creditAccount.EstablishmentID,
		ClientID:                creditAccount.ClientID,
		CreditLimit:             creditAccount.CreditLimit,
		CurrentBalance:          creditAccount.CurrentBalance,
		MonthlyDueDate:          creditAccount.MonthlyDueDate,
		InterestRate:            creditAccount.InterestRate,
		InterestType:            creditAccount.InterestType,
		CreditType:              creditAccount.CreditType,
		GracePeriod:             creditAccount.GracePeriod,
		IsBlocked:               creditAccount.IsBlocked,
		LastInterestAccrualDate: creditAccount.LastInterestAccrualDate,
		CreatedAt:               creditAccount.CreatedAt,
		UpdatedAt:               creditAccount.UpdatedAt,
		LateFeeRuleID:           creditAccount.LateFeeRuleID,
		Client:                  clientResponse,
	}
}

func NewClientResponse(client *entities.Client) *response.ClientResponse {
	if client == nil {
		return nil
	}
	userResponse := NewUserResponse(client.User)
	return &response.ClientResponse{
		ID:        client.ID,
		User:      userResponse,
		IsActive:  client.IsActive,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

func NewUserResponse(user *entities.User) *response.UserResponse {
	if user == nil {
		return nil
	}
	return &response.UserResponse{
		ID:        user.ID,
		DNI:       user.DNI,
		Name:      user.Name,
		Email:     user.Email,
		Address:   user.Address,
		Phone:     user.Phone,
		Rol:       user.Rol,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func responseToCreditAccount(res *response.CreditAccountResponse) *entities.CreditAccount {
	return &entities.CreditAccount{
		Model: gorm.Model{
			ID:        res.ID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
		EstablishmentID:         res.EstablishmentID,
		ClientID:                res.ClientID,
		CreditLimit:             res.CreditLimit,
		CurrentBalance:          res.CurrentBalance,
		MonthlyDueDate:          res.MonthlyDueDate,
		InterestRate:            res.InterestRate,
		InterestType:            res.InterestType,
		CreditType:              res.CreditType,
		GracePeriod:             res.GracePeriod,
		IsBlocked:               res.IsBlocked,
		LastInterestAccrualDate: res.LastInterestAccrualDate,
		LateFeeRuleID:           res.LateFeeRuleID,
		Client:                  responseToClient(res.Client),
		LateFeeRule:             responseToLateFeeRule(res.LateFeeRule),
	}
}

func responseToLateFeeRule(res *response.LateFeeRuleResponse) *entities.LateFeeRule {
	if res == nil {
		return nil
	}

	return &entities.LateFeeRule{
		Model: gorm.Model{
			ID: res.ID,
		},
		EstablishmentID: res.EstablishmentID,
		Name:            res.Name,
		DaysOverdueMin:  res.DaysOverdueMin,
		DaysOverdueMax:  res.DaysOverdueMax,
		FeeType:         res.FeeType,
		FeeValue:        res.FeeValue,
	}
}

func responseToClient(res *response.ClientResponse) *entities.Client {
	if res == nil {
		return nil
	}

	return &entities.Client{
		Model: gorm.Model{
			ID:        res.ID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
		IsActive: res.IsActive,
	}
}


func (s *creditAccountService) GetClientAccountStatement(clientID uint, startDate, endDate time.Time) ([]*response.AccountStatementResponse, error) {
    // 1. Get the client's credit account
    creditAccounts, err := s.creditAccountRepo.GetByClientID(clientID)
    if err != nil {
        return nil, fmt.Errorf("error retrieving credit account: %w", err)
    }
    if creditAccounts == nil {
        return nil, fmt.Errorf("credit account for client %d not found", clientID)
    }

   
    var statements []*response.AccountStatementResponse // Array to hold statements
    for _, creditAccount := range creditAccounts { // Iterate through credit accounts
        // 2. Get transactions within the date range
        transactions, err := s.transactionRepo.GetTransactionsByCreditAccountIDAndDateRange(creditAccount.ID, startDate, endDate)
        if err != nil {
            return nil, fmt.Errorf("error retrieving transactions: %w", err)
        }

        // 3. Calculate balance at the beginning of the statement period
        var startingBalance float64
        if !startDate.IsZero() {
            startingBalance, err = s.transactionRepo.GetBalanceBeforeDate(creditAccount.ID, startDate)
            if err != nil {
                return nil, fmt.Errorf("error calculating starting balance: %w", err)
            }
        } else {
            startingBalance = 0.0 
        }

        // 4. Create AccountStatementResponse
        statement := &response.AccountStatementResponse{
            ClientID:        clientID,
            StartDate:       startDate,
            EndDate:         endDate,
            StartingBalance: startingBalance,
            Transactions:    make([]response.TransactionResponse, len(transactions)),
        }

       	// 5. Populate transactions in the statement
		for i, transaction := range transactions {
    		statement.Transactions[i] = response.TransactionResponse {
        		ID:               transaction.ID,
        		CreditAccountID:  transaction.CreditAccountID,
        		TransactionType:  transaction.TransactionType, 
        		Amount:           transaction.Amount,
        		Description:      transaction.Description,
        		CreatedAt:        transaction.CreatedAt,
    		}
		}

        statements = append(statements, statement)
    }

    return statements, nil
}

func (s *creditAccountService) CalculateDueDate(account entities.CreditAccount) (time.Time, error) {
    today := time.Now()

    if account.CreditType == enums.ShortTerm {
        // Short-term credit: Due date is the next month's due date
        nextMonth := today.Month() + 1
        nextYear := today.Year()
        if nextMonth > time.December {
            nextMonth = time.January
            nextYear++
        }
        return time.Date(nextYear, nextMonth, account.MonthlyDueDate, 0, 0, 0, 0, time.UTC), nil
    } else if account.CreditType == enums.LongTerm {
        // Long-term credit: Due date depends on installment schedule

        // 1. Retrieve installments for the credit account
        installments, err := s.installmentRepo.GetByCreditAccountID(account.ID) 
        if err != nil {
            return time.Time{}, fmt.Errorf("error retrieving installments: %w", err)
        }

        // 2. Find the next pending installment
        var nextDueDate time.Time
        for _, installment := range installments {
            if installment.Status == enums.Pending && installment.DueDate.After(today) {
                nextDueDate = installment.DueDate
                break
            }
        }

        // 3. If no pending installments, calculate next due date based on MonthlyDueDate
        if nextDueDate.IsZero() {
            nextMonth := today.Month() + 1
            nextYear := today.Year()
            if nextMonth > time.December {
                nextMonth = time.January
                nextYear++
            }
            nextDueDate = time.Date(nextYear, nextMonth, account.MonthlyDueDate, 0, 0, 0, 0, time.UTC)
        }

        return nextDueDate, nil 
    }

    // Handle invalid credit type (consider adding error logging)
    return time.Time{}, fmt.Errorf("invalid credit type: %s", account.CreditType)
}

// GetClientAccountHistory retrieves the complete transaction history for a client's credit account.
func (s *creditAccountService) GetClientAccountHistory(clientID uint) (*response.AccountStatementResponse, error) {
    // Get the client's credit account
    creditAccounts, err := s.creditAccountRepo.GetByClientID(clientID)
    if err != nil {
        return nil, fmt.Errorf("error retrieving credit account: %w", err)
    }
    if creditAccounts == nil {
        return nil, fmt.Errorf("credit account for client %d not found", clientID)
    }

    var allTransactions []entities.Transaction
    var startingBalance float64
    for _, creditAccount := range creditAccounts {
        // 1. Get all transactions for the credit account (no date filter)
        transactions, err := s.transactionRepo.GetTransactionsByCreditAccountID(creditAccount.ID)
        if err != nil {
            return nil, fmt.Errorf("error retrieving transactions: %w", err)
        }
        allTransactions = append(allTransactions, transactions...)

        // 2. Update Starting Balance
        balance, err := s.creditAccountRepo.GetStartingBalance(creditAccount.ID)
        if err != nil {
            return nil, fmt.Errorf("error getting starting balance: %w", err)
        }
        startingBalance += balance
    }

    // 3. Create AccountStatementResponse (start and end dates are empty)
    statement := &response.AccountStatementResponse{
        ClientID:        clientID,
        StartingBalance: startingBalance,
        Transactions:    make([]response.TransactionResponse, len(allTransactions)),
    }

    // 4. Populate transactions in the statement
    for i, transaction := range allTransactions {
        statement.Transactions[i] = response.TransactionResponse{
            ID:                transaction.ID,
            CreditAccountID:  transaction.CreditAccountID,
            TransactionType:   transaction.TransactionType,
            Amount:            transaction.Amount,
            Description:       transaction.Description,
            CreatedAt:         transaction.CreatedAt,
        }
    }

    return statement, nil
}