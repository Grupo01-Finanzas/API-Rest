package service

import (
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"errors"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"math"
	"os"
	"time"
)

// PurchaseService handles purchase logic.
type PurchaseService interface {
	ProcessPurchase(userID uint, establishmentID uint, productIDs []uint, creditType enums.CreditType, amount float64) error
	GetClientBalance(clientID uint) (float64, error)
	GetClientOverdueBalance(clientID uint) (float64, error)
	GetClientInstallments(clientID uint) ([]response.InstallmentResponse, error)
	GetClientTransactions(clientID uint) ([]response.TransactionResponse, error)
	GetClientCreditAccount(clientID uint) (*entities.CreditAccount, error)
	GetClientAccountSummary(clientID uint) (*response.AccountSummaryResponse, error)
	CalculateDueDate(account entities.CreditAccount) (time.Time, error)
	GetClientAccountStatement(clientID uint, startDate, endDate time.Time) (*response.AccountStatementResponse, error)
	GenerateClientAccountStatementPDF(clientID uint, startDate, endDate time.Time) ([]byte, error)
}

type purchaseService struct {
	userRepo          repository.UserRepository
	establishmentRepo repository.EstablishmentRepository
	productRepo       repository.ProductRepository
	creditAccountRepo repository.CreditAccountRepository
	transactionRepo   repository.TransactionRepository
	installmentRepo   repository.InstallmentRepository
}

func NewPurchaseService(userRepo repository.UserRepository, establishmentRepo repository.EstablishmentRepository, productRepo repository.ProductRepository, creditAccountRepo repository.CreditAccountRepository, transactionRepo repository.TransactionRepository, installmentRepo repository.InstallmentRepository) PurchaseService {
	return &purchaseService{
		userRepo:          userRepo,
		establishmentRepo: establishmentRepo,
		productRepo:       productRepo,
		creditAccountRepo: creditAccountRepo,
		transactionRepo:   transactionRepo,
		installmentRepo:   installmentRepo,
	}
}

func (s *purchaseService) ProcessPurchase(userID uint, establishmentID uint, productIDs []uint, creditType enums.CreditType, amount float64) error {
	if userID == 0 || establishmentID == 0 || len(productIDs) == 0 || amount <= 0 {
		return errors.New("invalid input data")
	}

	if creditType != enums.ShortTerm && creditType != enums.LongTerm {
		return errors.New("invalid credit type")
	}

	// Get the client's credit account
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(userID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return errors.New("client does not have a credit account")
	}

	// Check if the account is blocked
	if creditAccount.IsBlocked {
		return errors.New("client's credit account is blocked")
	}

	// Check if the purchase exceeds the credit limit
	if creditAccount.CurrentBalance+amount > creditAccount.CreditLimit {
		return fmt.Errorf("purchase amount exceeds credit limit (Current Balance: %.2f, Credit Limit: %.2f)", creditAccount.CurrentBalance, creditAccount.CreditLimit)
	}

	// If long-term credit, calculate and create installments
	if creditType == enums.LongTerm {
		err = s.createInstallments(creditAccount, amount)
		if err != nil {
			return fmt.Errorf("error creating installments: %w", err)
		}
	}

	// Start a transaction to ensure data consistency
	if err := s.creditAccountRepo.ProcessPurchaseTransaction(creditAccount, amount, "Product Purchase"); err != nil {
		return fmt.Errorf("error processing purchase: %w", err)
	}

	return nil

}

func (s *purchaseService) GetClientBalance(clientID uint) (float64, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return 0, errors.New("client does not have a credit account")
	}
	return creditAccount.CurrentBalance, nil
}

func (s *purchaseService) GetClientOverdueBalance(clientID uint) (float64, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return 0, nil // No credit account, no overdue balance
	}

	if !isAccountOverdue(*creditAccount) {
		return 0, nil // Account is not overdue
	}

	return creditAccount.CurrentBalance, nil
}

// isAccountOverdue checks if the account is overdue based on the monthly due date
func isAccountOverdue(creditAccount entities.CreditAccount) bool {
	today := time.Now()
	dueDate := time.Date(today.Year(), today.Month(), creditAccount.MonthlyDueDate, 0, 0, 0, 0, time.UTC)
	return today.After(dueDate) && creditAccount.CurrentBalance > 0
}

func (s *purchaseService) GetClientInstallments(clientID uint) ([]response.InstallmentResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return nil, errors.New("client does not have a credit account")
	}

	installments, err := s.installmentRepo.GetInstallmentsByCreditAccountID(creditAccount.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving installments: %w", err)
	}

	var installmentResponses []response.InstallmentResponse
	for _, installment := range installments {
		installmentResponses = append(installmentResponses, response.InstallmentResponse{
			ID:              installment.ID,
			CreditAccountID: installment.CreditAccountID,
			DueDate:         installment.DueDate,
			Amount:          installment.Amount,
			Status:          installment.Status,
			CreatedAt:       installment.CreatedAt,
			UpdatedAt:       installment.UpdatedAt,
		})
	}
	return installmentResponses, nil
}

func (s *purchaseService) GetClientTransactions(clientID uint) ([]response.TransactionResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return nil, errors.New("client does not have a credit account")
	}

	transactions, err := s.transactionRepo.GetTransactionsByCreditAccountID(creditAccount.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions: %w", err)
	}

	var transactionResponses []response.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, response.TransactionResponse{
			ID:              transaction.ID,
			CreditAccountID: transaction.CreditAccountID,
			TransactionType: transaction.TransactionType,
			Amount:          transaction.Amount,
			Description:     transaction.Description,
			TransactionDate: transaction.TransactionDate,
			CreatedAt:       transaction.CreatedAt,
			UpdatedAt:       transaction.UpdatedAt,
		})
	}

	return transactionResponses, nil
}

func (s *purchaseService) GetClientCreditAccount(clientID uint) (*entities.CreditAccount, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByClientID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return nil, errors.New("client does not have a credit account")
	}
	return creditAccount, nil
}

func (s *purchaseService) createInstallments(creditAccount *entities.CreditAccount, purchaseAmount float64) error {
	if creditAccount.CreditType != enums.LongTerm {
		return nil // Installments are not applicable for short-term credit
	}

	// Assuming 12-month installment plan for simplicity
	numInstallments := 12
	installmentAmount := purchaseAmount / float64(numInstallments)

	// Calculate the first installment due date based on credit account's due date
	firstDueDate := calculateNextDueDate(creditAccount.MonthlyDueDate)

	var installments []entities.Installment
	for i := 0; i < numInstallments; i++ {
		installmentDueDate := firstDueDate.AddDate(0, i, 0)
		installment := entities.Installment{
			CreditAccountID: creditAccount.ID,
			DueDate:         installmentDueDate,
			Amount:          installmentAmount,
			Status:          enums.Pending,
		}
		installments = append(installments, installment)
	}

	return s.installmentRepo.CreateInstallments(installments)
}

// calculateNextDueDate calculates the next due date for an installment
func calculateNextDueDate(monthlyDueDate int) time.Time {
	today := time.Now()
	dueDate := time.Date(today.Year(), today.Month(), monthlyDueDate, 0, 0, 0, 0, time.UTC)
	if dueDate.Before(today) {
		dueDate = dueDate.AddDate(0, 1, 0)
	}
	return dueDate
}

// CalculateDueDate calculates the next due date for a credit account.
func (s *purchaseService) CalculateDueDate(account entities.CreditAccount) (time.Time, error) {
	today := time.Now()
	if account.CreditType == enums.ShortTerm {
		// For short-term credit, the due date is the next month's due date
		nextMonth := today.Month() + 1
		nextYear := today.Year()
		if nextMonth > time.December {
			nextMonth = time.January
			nextYear++
		}
		return time.Date(nextYear, nextMonth, account.MonthlyDueDate, 0, 0, 0, 0, time.UTC), nil
	} else if account.CreditType == enums.LongTerm {
		// For long-term credit, find the next pending installment's due date
		installments, err := s.installmentRepo.GetInstallmentsByCreditAccountID(account.ID)
		if err != nil {
			return time.Time{}, fmt.Errorf("error retrieving installments: %w", err)
		}
		for _, installment := range installments {
			if installment.Status == enums.Pending && installment.DueDate.After(today) {
				return installment.DueDate, nil
			}
		}

		// If no pending installments, calculate the next due date based on MonthlyDueDate
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

// GetClientAccountSummary retrieves a summary of the client's account.
func (s *purchaseService) GetClientAccountSummary(clientID uint) (*response.AccountSummaryResponse, error) {
	creditAccount, err := s.GetClientCreditAccount(clientID)
	if err != nil {
		return nil, err
	}

	// Get transactions up to the current due date
	dueDate, err := s.CalculateDueDate(*creditAccount)
	if err != nil {
		return nil, fmt.Errorf("error calculating due date: %w", err)
	}

	transactions, err := s.transactionRepo.GetTransactionsByCreditAccountIDAndDateRange(creditAccount.ID, time.Time{}, dueDate)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions: %w", err)
	}

	// Calculate total interest accrued on purchases (you'll need to implement calculateInterestForTransactions)
	totalInterest := calculateInterestForTransactions(transactions, *creditAccount, dueDate)

	// Prepare the response
	summary := &response.AccountSummaryResponse{
		CurrentBalance: creditAccount.CurrentBalance,
		DueDate:        dueDate,
		TotalInterest:  totalInterest,
		Transactions:   make([]response.TransactionResponse, len(transactions)),
	}

	// Populate transactions in the response
	for i, transaction := range transactions {
		summary.Transactions[i] = *transactionToResponse(&transaction)
	}

	return summary, nil
}

// calculateInterestForTransactions calculates interest for a list of transactions.
func calculateInterestForTransactions(transactions []entities.Transaction, account entities.CreditAccount, dueDate time.Time) float64 {
	var totalInterest float64

	for _, transaction := range transactions {
		if transaction.TransactionType == enums.Purchase {
			interest := calculateInterestForPurchase(transaction, account, dueDate)
			totalInterest += interest
		}
	}

	return totalInterest
}

// calculateInterestForPurchase calculates interest for a single purchase transaction.
func calculateInterestForPurchase(transaction entities.Transaction, account entities.CreditAccount, dueDate time.Time) float64 {
	// Calculate the number of days from the purchase date to the due date
	days := int(dueDate.Sub(transaction.TransactionDate).Hours() / 24)

	// Calculate the interest based on the interest type (Nominal or Effective)
	var interest float64
	principal := transaction.Amount
	annualRate := account.InterestRate / 100

	if account.InterestType == enums.Nominal {
		// Nominal interest calculation
		interest = principal * annualRate * float64(days) / 365
	} else if account.InterestType == enums.Effective {
		// Effective interest calculation
		dailyRate := math.Pow(1+annualRate, 1.0/365) - 1
		interest = principal * (math.Pow(1+dailyRate, float64(days)) - 1)
	}

	return interest
}

// GetClientAccountStatement retrieves a client's account statement for a date range.
func (s *purchaseService) GetClientAccountStatement(clientID uint, startDate, endDate time.Time) (*response.AccountStatementResponse, error) {
	creditAccount, err := s.GetClientCreditAccount(clientID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transactionRepo.GetTransactionsByCreditAccountIDAndDateRange(creditAccount.ID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions: %w", err)
	}

	// Calculate starting balance (you might need a helper function in your repository)
	startingBalance := 0.0
	if !startDate.IsZero() {
		startingBalance, err = s.transactionRepo.GetBalanceBeforeDate(creditAccount.ID, startDate)
		if err != nil {
			return nil, fmt.Errorf("error getting starting balance: %w", err)
		}
	}

	// Prepare the response
	statement := &response.AccountStatementResponse{
		ClientID:        clientID,
		StartDate:       startDate,
		EndDate:         endDate,
		StartingBalance: startingBalance,
		Transactions:    make([]response.TransactionResponse, len(transactions)),
	}

	// Populate transactions in the response
	for i, transaction := range transactions {
		statement.Transactions[i] = *transactionToResponse(&transaction)
	}

	return statement, nil
}

// GenerateClientAccountStatementPDF generates a PDF account statement for the client.
func (s *purchaseService) GenerateClientAccountStatementPDF(clientID uint, startDate, endDate time.Time) ([]byte, error) {
	// 1. Get account statement data
	statement, err := s.GetClientAccountStatement(clientID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting account statement: %w", err)
	}

	// 2. Generate PDF using the statement data
	pdf := gofpdf.New("P", "mm", "A4", "") // Create a new PDF document
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Account Statement - Client ID: %d", clientID))
	pdf.Ln(10)

	// Date Range
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 10, fmt.Sprintf("Start Date: %s", startDate.Format("2006-01-02")), "", 0, "L", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("End Date: %s", endDate.Format("2006-01-02")), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Starting Balance
	pdf.CellFormat(40, 10, fmt.Sprintf("Starting Balance: %.2f", statement.StartingBalance), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Transactions Table Header (Corrected)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(30, 10, "Date")
	pdf.Cell(40, 10, "Description")
	pdf.Cell(30, 10, "Type")
	pdf.Cell(30, 10, "Payment Method")
	pdf.Cell(30, 10, "Amount")
	pdf.Cell(30, 10, "Status")
	pdf.Ln(10)

	// Transactions Table Data
	pdf.SetFont("Arial", "", 10)
	for _, transaction := range statement.Transactions {
		pdf.CellFormat(30, 10, transaction.TransactionDate.Format("2006-01-02"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, transaction.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, string(transaction.TransactionType), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, string(transaction.PaymentMethod), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", transaction.Amount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 10, string(transaction.PaymentStatus), "1", 0, "L", false, 0, "")
		pdf.Ln(8)
	}

	// Ending Balance
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, fmt.Sprintf("Ending Balance: %.2f", statement.StartingBalance+calculateTotalTransactionAmount(statement.Transactions)), "", 0, "L", false, 0, "")

	// 3. Output PDF as byte array
	err = pdf.OutputFileAndClose("account_statement.pdf") // Correct way to output to file
	if err != nil {
		return nil, fmt.Errorf("error outputting PDF to file: %w", err)
	}

	// Read the PDF file contents into a byte array
	pdfBytes, err := os.ReadFile("account_statement.pdf")
	if err != nil {
		return nil, fmt.Errorf("error reading PDF file: %w", err)
	}

	// Optionally, you can delete the file after reading it:
	// os.Remove("account_statement.pdf")

	return pdfBytes, nil
}

// calculateTotalTransactionAmount calculates the total amount from a list of transactions
func calculateTotalTransactionAmount(transactions []response.TransactionResponse) float64 {
	total := 0.0
	for _, transaction := range transactions {
		if transaction.TransactionType == enums.Purchase {
			total += transaction.Amount
		} else if transaction.TransactionType == enums.Payment {
			total -= transaction.Amount
		}
	}
	return total
}
