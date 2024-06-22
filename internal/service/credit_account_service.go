package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"errors"
	"fmt"
	"time"
)

// CreditAccountService handles credit account-related operations.
type CreditAccountService interface {
	CreateCreditAccount(req request.CreateCreditAccountRequest, establishmentID uint) (*response.CreditAccountResponse, error)
	GetCreditAccountByID(id uint) (*response.CreditAccountResponse, error)
	UpdateCreditAccount(id uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error)
	DeleteCreditAccount(id uint) error
	GetCreditAccountsByEstablishmentID(establishmentID uint) ([]response.CreditAccountResponse, error)
	GetCreditAccountByClientID(clientID uint) (*response.CreditAccountResponse, error)
	ApplyInterestToAccount(creditAccountID uint) error
	ApplyLateFeeToAccount(creditAccountID uint) error
	GetOverdueCreditAccounts(establishmentID uint) ([]response.CreditAccountResponse, error)
	ProcessPurchase(creditAccountID uint, amount float64, description string) error
	ProcessPayment(creditAccountID uint, amount float64, description string) error
	GetAdminDebtSummary(establishmentID uint) ([]response.AdminDebtSummary, error)
	CalculateDueDate(account entities.CreditAccount) (time.Time, error)
	GetNumberOfDues(account entities.CreditAccount) int
	UpdateCreditAccountByClientID(clientID uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error)
}

type creditAccountService struct {
	creditAccountRepo repository.CreditAccountRepository
	transactionRepo   repository.TransactionRepository
	installmentRepo   repository.InstallmentRepository
	clientRepo        repository.ClientRepository
	establishmentRepo repository.EstablishmentRepository
}

// NewCreditAccountService creates a new instance of CreditAccountService.
func NewCreditAccountService(creditAccountRepo repository.CreditAccountRepository, transactionRepo repository.TransactionRepository, installmentRepo repository.InstallmentRepository, clientRepo repository.ClientRepository, establishmentRepo repository.EstablishmentRepository) CreditAccountService {
	return &creditAccountService{
		creditAccountRepo: creditAccountRepo,
		transactionRepo:   transactionRepo,
		installmentRepo:   installmentRepo,
		clientRepo:        clientRepo,
		establishmentRepo: establishmentRepo,
	}
}

// CreateCreditAccount creates a new credit account for a client.
func (s *creditAccountService) CreateCreditAccount(req request.CreateCreditAccountRequest, establishmentID uint) (*response.CreditAccountResponse, error) {
	client, err := s.clientRepo.GetClientByID(req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving client: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("client with ID %d not found", req.ClientID)
	}

	// Check if the establishment exists
	establishment, err := s.establishmentRepo.GetEstablishmentByID(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving establishment: %w", err)
	}
	if establishment == nil {
		return nil, fmt.Errorf("establishment with ID %d not found", establishmentID)
	}

	creditAccount := entities.CreditAccount{
		EstablishmentID:         establishment.ID,
		ClientID:                client.ID,
		CreditLimit:             req.CreditLimit,
		MonthlyDueDate:          req.MonthlyDueDate,
		InterestRate:            req.InterestRate,
		InterestType:            req.InterestType,
		CreditType:              req.CreditType,
		GracePeriod:             req.GracePeriod,
		IsBlocked:               false,
		LastInterestAccrualDate: time.Now(),
		CurrentBalance:          0,
		LateFeePercentage:       establishment.LateFeePercentage, // Use the establishment's late fee percentage
	}

	err = s.creditAccountRepo.CreateCreditAccount(&creditAccount)
	if err != nil {
		return nil, fmt.Errorf("error creating credit account: %w", err)
	}

	return creditAccountToResponse(&creditAccount), nil
}

// GetCreditAccountByID retrieves a credit account by ID.
func (s *creditAccountService) GetCreditAccountByID(id uint) (*response.CreditAccountResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(id)
	if err != nil {
		return nil, err
	}

	return creditAccountToResponse(creditAccount), nil
}

// UpdateCreditAccount updates an existing credit account.
func (s *creditAccountService) UpdateCreditAccount(id uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields only if they are provided in the request
	if req.CreditLimit > 0 {
		creditAccount.CreditLimit = req.CreditLimit
	}
	if req.MonthlyDueDate > 0 {
		creditAccount.MonthlyDueDate = req.MonthlyDueDate
	}
	if req.InterestRate > 0 {
		creditAccount.InterestRate = req.InterestRate
	}
	if req.InterestType != "" {
		creditAccount.InterestType = req.InterestType
	}
	if req.CreditType != "" {
		creditAccount.CreditType = req.CreditType
	}
	if req.GracePeriod >= 0 {
		creditAccount.GracePeriod = req.GracePeriod
	}
	creditAccount.IsBlocked = req.IsBlocked
	if req.LateFeePercentage >= 0 {
		creditAccount.LateFeePercentage = req.LateFeePercentage
	}

	err = s.creditAccountRepo.UpdateCreditAccount(creditAccount)
	if err != nil {
		return nil, err
	}

	return creditAccountToResponse(creditAccount), nil
}

// DeleteCreditAccount deletes a credit account.
func (s *creditAccountService) DeleteCreditAccount(id uint) error {
	return s.creditAccountRepo.DeleteCreditAccount(id)
}

// GetCreditAccountsByEstablishmentID retrieves all credit accounts for an establishment.
func (s *creditAccountService) GetCreditAccountsByEstablishmentID(establishmentID uint) ([]response.CreditAccountResponse, error) {
	creditAccounts, err := s.creditAccountRepo.GetCreditAccountsByEstablishmentID(establishmentID)
	if err != nil {
		return nil, err
	}

	var creditAccountResponses []response.CreditAccountResponse
	for _, creditAccount := range creditAccounts {
		creditAccountResponses = append(creditAccountResponses, *creditAccountToResponse(&creditAccount))
	}
	return creditAccountResponses, nil
}

// GetCreditAccountByClientID retrieves a credit account by client ID.
func (s *creditAccountService) GetCreditAccountByClientID(clientID uint) (*response.CreditAccountResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	return creditAccountToResponse(creditAccount), nil
}

// ApplyInterestToAccount calculates and applies interest to a credit account.
func (s *creditAccountService) ApplyInterestToAccount(creditAccountID uint) error {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(creditAccountID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}

	if err := s.creditAccountRepo.ApplyInterest(creditAccount); err != nil {
		return fmt.Errorf("error applying interest to account %d: %w", creditAccountID, err)
	}
	return nil
}

// ApplyLateFeeToAccount applies late fee to a credit account if overdue.
func (s *creditAccountService) ApplyLateFeeToAccount(creditAccountID uint) error {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(creditAccountID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}

	// Calculate days overdue (you can use a helper function for this)
	daysOverdue := calculateDaysOverdue(creditAccount.MonthlyDueDate)

	if err := s.creditAccountRepo.ApplyLateFee(creditAccount, daysOverdue); err != nil {
		return fmt.Errorf("error applying late fee to account %d: %w", creditAccountID, err)
	}
	return nil
}

// calculateDaysOverdue calculates the number of days a payment is overdue
func calculateDaysOverdue(dueDate int) int {
	today := time.Now()
	thisMonthDueDate := time.Date(today.Year(), today.Month(), dueDate, 0, 0, 0, 0, time.UTC)

	if today.Before(thisMonthDueDate) {
		return 0
	}

	return int(today.Sub(thisMonthDueDate).Hours() / 24)
}

// GetOverdueCreditAccounts retrieves overdue credit accounts for an establishment.
func (s *creditAccountService) GetOverdueCreditAccounts(establishmentID uint) ([]response.CreditAccountResponse, error) {
	overdueAccounts, err := s.creditAccountRepo.GetOverdueCreditAccounts(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving overdue credit accounts: %w", err)
	}

	var overdueAccountResponses []response.CreditAccountResponse
	for _, account := range overdueAccounts {
		overdueAccountResponses = append(overdueAccountResponses, *creditAccountToResponse(&account))
	}

	return overdueAccountResponses, nil
}

// ProcessPurchase processes a purchase transaction on a credit account.
func (s *creditAccountService) ProcessPurchase(creditAccountID uint, amount float64, description string) error {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(creditAccountID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}

	return s.creditAccountRepo.ProcessPurchase(creditAccount, amount, description)
}

// ProcessPayment processes a payment transaction on a credit account.
func (s *creditAccountService) ProcessPayment(creditAccountID uint, amount float64, description string) error {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(creditAccountID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}

	return s.creditAccountRepo.ProcessPayment(creditAccount, amount, description)
}

// GetAdminDebtSummary retrieves a summary of debts for an establishment.
func (s *creditAccountService) GetAdminDebtSummary(establishmentID uint) ([]response.AdminDebtSummary, error) {
	creditAccounts, err := s.creditAccountRepo.GetCreditAccountsByEstablishmentID(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit accounts: %w", err)
	}

	summary := make([]response.AdminDebtSummary, 0, len(creditAccounts))
	for _, account := range creditAccounts {
		client, err := s.clientRepo.GetClientByID(account.ClientID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving client: %w", err)
		}
		user := client.User

		dueDate, err := s.CalculateDueDate(account)
		if err != nil {
			return nil, fmt.Errorf("error calculating due date: %w", err)
		}

		summaryItem := response.AdminDebtSummary{
			ClientID:       account.ClientID,
			ClientName:     user.Name,
			CreditType:     string(account.CreditType),
			InterestRate:   account.InterestRate,
			NumberOfDues:   s.GetNumberOfDues(account),
			CurrentBalance: account.CurrentBalance,
			DueDate:        dueDate,
		}

		summary = append(summary, summaryItem)
	}
	return summary, nil
}

// CalculateDueDate calculates the next due date for a credit account.
func (s *creditAccountService) CalculateDueDate(account entities.CreditAccount) (time.Time, error) {
	today := time.Now()
	if account.CreditType == enums.ShortTerm {
		nextMonth := today.Month() + 1
		nextYear := today.Year()
		if nextMonth > time.December {
			nextMonth = time.January
			nextYear++
		}
		return time.Date(nextYear, nextMonth, account.MonthlyDueDate, 0, 0, 0, 0, time.UTC), nil
	} else if account.CreditType == enums.LongTerm {
		installments, err := s.installmentRepo.GetInstallmentsByCreditAccountID(account.ID)
		if err != nil {
			return time.Time{}, fmt.Errorf("error retrieving installments: %w", err)
		}
		for _, installment := range installments {
			if installment.Status == enums.Pending && installment.DueDate.After(today) {
				return installment.DueDate, nil
			}
		}
		nextMonth := today.Month() + 1
		nextYear := today.Year()
		if nextMonth > time.December {
			nextMonth = time.January
			nextYear++
		}
		return time.Date(nextYear, nextMonth, account.MonthlyDueDate, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("invalid credit type: %s", account.CreditType)
}

func (s *creditAccountService) GetNumberOfDues(account entities.CreditAccount) int {
	if account.CreditType != enums.LongTerm {
		return 0
	}

	installments, err := s.installmentRepo.GetInstallmentsByCreditAccountID(account.ID)
	if err != nil {
		return 0
	}
	return len(installments)
}

func creditAccountToResponse(creditAccount *entities.CreditAccount) *response.CreditAccountResponse {
	return &response.CreditAccountResponse{
		ID:                      creditAccount.ID,
		ClientID:                creditAccount.ClientID,
		Client:                  NewUserResponse(creditAccount.Client),
		EstablishmentID:         creditAccount.EstablishmentID,
		Establishment:           _NewEstablishmentResponse(creditAccount.Establishment),
		CreditLimit:             creditAccount.CreditLimit,
		CurrentBalance:          creditAccount.CurrentBalance,
		MonthlyDueDate:          creditAccount.MonthlyDueDate,
		InterestRate:            creditAccount.InterestRate,
		InterestType:            creditAccount.InterestType,
		CreditType:              creditAccount.CreditType,
		GracePeriod:             creditAccount.GracePeriod,
		IsBlocked:               creditAccount.IsBlocked,
		LastInterestAccrualDate: creditAccount.LastInterestAccrualDate,
		LateFeePercentage:       creditAccount.LateFeePercentage,
		CreatedAt:               creditAccount.CreatedAt,
		UpdatedAt:               creditAccount.UpdatedAt,
	}
}

func _NewEstablishmentResponse(establishment *entities.Establishment) *response.EstablishmentResponse {
	if establishment == nil {
		return nil
	}
	return &response.EstablishmentResponse{
		ID:                establishment.ID,
		RUC:               establishment.RUC,
		Name:              establishment.Name,
		Phone:             establishment.Phone,
		Address:           establishment.Address,
		ImageUrl:          establishment.ImageUrl,
		LateFeePercentage: establishment.LateFeePercentage,
		IsActive:          establishment.IsActive,
		CreatedAt:         establishment.CreatedAt,
		UpdatedAt:         establishment.UpdatedAt,
	}
}

// UpdateCreditAccountByClientID updates an existing credit account by client ID.
func (s *creditAccountService) UpdateCreditAccountByClientID(clientID uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return nil, errors.New("credit account not found for this client")
	}

	// Update the credit account fields based on the request
	if req.CreditLimit > 0 {
		creditAccount.CreditLimit = req.CreditLimit
	}
	if req.MonthlyDueDate > 0 {
		creditAccount.MonthlyDueDate = req.MonthlyDueDate
	}
	if req.InterestRate > 0 {
		creditAccount.InterestRate = req.InterestRate
	}
	if req.InterestType != "" {
		creditAccount.InterestType = req.InterestType
	}
	if req.CreditType != "" {
		creditAccount.CreditType = req.CreditType
	}
	if req.GracePeriod >= 0 {
		creditAccount.GracePeriod = req.GracePeriod
	}
	creditAccount.IsBlocked = req.IsBlocked
	if req.LateFeePercentage >= 0 {
		creditAccount.LateFeePercentage = req.LateFeePercentage
	}

	err = s.creditAccountRepo.UpdateCreditAccount(creditAccount)
	if err != nil {
		return nil, fmt.Errorf("error updating credit account: %w", err)
	}

	return creditAccountToResponse(creditAccount), nil
}
