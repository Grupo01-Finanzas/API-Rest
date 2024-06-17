package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"gorm.io/gorm"
)

// PurchaseService handles purchase logic
type PurchaseService interface {
	ProcessPurchase(userID uint, establishmentID uint, productIDs []uint, creditType enums.CreditType, amount float64) error
	GetClientBalance(clientID uint) (float64, error)
	GetClientOverdueBalance(clientID uint) (float64, error)
	GetClientInstallments(clientID uint) ([]response.InstallmentResponse, error)
	GetClientTransactions(clientID uint) ([]response.TransactionResponse, error)
	GetClientCreditAccount(clientID uint) (*entities.CreditAccount, error)
}

type purchaseService struct {
	userRepo          repository.UserRepository
	establishmentRepo repository.EstablishmentRepository
	productRepo       repository.ProductRepository
	creditAccountRepo repository.CreditAccountRepository
	transactionRepo   repository.TransactionRepository
	installmentRepo   repository.InstallmentRepository
	db                *gorm.DB
}

func NewPurchaseService(userRepo repository.UserRepository, establishmentRepo repository.EstablishmentRepository, productRepo repository.ProductRepository, creditAccountRepo repository.CreditAccountRepository, transactionRepo repository.TransactionRepository, installmentRepo repository.InstallmentRepository, db *gorm.DB) PurchaseService {
	return &purchaseService{
		userRepo:          userRepo,
		establishmentRepo: establishmentRepo,
		productRepo:       productRepo,
		creditAccountRepo: creditAccountRepo,
		transactionRepo:   transactionRepo,
		installmentRepo:   installmentRepo,
		db:                db,
	}
}

func (s *purchaseService) ProcessPurchase(userID uint, establishmentID uint, productIDs []uint, creditType enums.CreditType, amount float64) error {
	// 1. Input data validation
	if userID == 0 || establishmentID == 0 || len(productIDs) == 0 || amount <= 0 {
		return errors.New("invalid input data")
	}
	if creditType != enums.ShortTerm && creditType != enums.LongTerm {
		return errors.New("invalid credit type")
	}

	// 2. Check if the user has an overdue balance
	overdueBalance, err := s.GetClientOverdueBalance(userID) // Using the service method
	if err != nil {
		return fmt.Errorf("error checking overdue balance: %w", err)
	}
	if overdueBalance > 0 {
		return errors.New("the user has an overdue balance, the purchase cannot be made")
	}

	// 3. Get or Create the user's credit account
	creditAccount, err := s.getOrCreateCreditAccount(userID, establishmentID, creditType)
	if err != nil {
		return fmt.Errorf("error getting/creating credit account: %w", err)
	}

	// 4. Verify credit limit
	if creditAccount.CurrentBalance+amount > creditAccount.CreditLimit {
		return errors.New("the purchase exceeds the user's credit limit")
	}

	// 5. Process the purchase within a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 5.1 Create purchase transaction
		transaction := &entities.Transaction{
			CreditAccountID: creditAccount.ID,
			TransactionType: enums.Purchase,
			Amount:          amount,
			Description:     "Product Purchase",
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("error creating transaction: %w", err)
		}

		// 5.2 Update credit account balance
		creditAccount.CurrentBalance += amount
		if err := tx.Save(&creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		// 5.3 Create installments if it's a long-term credit
		if creditType == enums.LongTerm {
			if err := s.createInstallments(tx, creditAccount, amount); err != nil {
				return fmt.Errorf("error creating installments: %w", err)
			}
		}

		// 5.4 Convert user to client if it's their first purchase
		if err := s.convertUserToClient(tx, userID); err != nil {
			return fmt.Errorf("error converting user to client: %w", err)
		}

		return nil // Transaction completed successfully
	})
}

// getOrCreateCreditAccount retrieves or creates the credit account
func (s *purchaseService) getOrCreateCreditAccount(userID uint, establishmentID uint, creditType enums.CreditType) (*entities.CreditAccount, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Get data to create the account (credit limit, monthly due date, etc.)
			// Here you could have logic to get these data from the request or configuration.
			creditLimit := 1000.00        // Example: Default credit limit
			monthlyDueDate := 15          // Example: Default monthly due date
			interestRate := 10.0          // Example: Default interest rate
			interestType := enums.Nominal // Example: Default interest type
			gracePeriod := 0              // Example: Default grace period

			// Get the late fee rule (you might have a default one)
			lateFeeRule, err := s.getLateFeeRule(establishmentID)
			if err != nil {
				return nil, fmt.Errorf("error getting late fee rule: %w", err)
			}

			creditAccount = &entities.CreditAccount{
				EstablishmentID:         establishmentID,
				ClientID:                userID,
				CreditLimit:             creditLimit,
				MonthlyDueDate:          monthlyDueDate,
				InterestRate:            interestRate,
				InterestType:            interestType,
				CreditType:              creditType,
				GracePeriod:             gracePeriod,
				IsBlocked:               false,
				LastInterestAccrualDate: time.Now(),
				CurrentBalance:          0.0,
				LateFeeRuleID:           lateFeeRule.ID,
			}

			if err := s.creditAccountRepo.Create(creditAccount); err != nil {
				return nil, fmt.Errorf("error creating credit account: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error retrieving credit account: %w", err)
		}
	}
	return creditAccount, nil
}

// getLateFeeRule gets the late fee rule for the establishment
func (s *purchaseService) getLateFeeRule(establishmentID uint) (*entities.LateFeeRule, error) {
	// ... Logic to retrieve the late fee rule for the establishment ...

	// Example: Get the first rule for the establishment
	rules, err := s.establishmentRepo.GetLateFeeRules(establishmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting late fee rules: %w", err)
	}
	if len(rules) == 0 {
		return nil, errors.New("no late fee rules found for the establishment")
	}

	return &rules[0], nil
}

// createInstallments creates the installments for the long-term credit
func (s *purchaseService) createInstallments(tx *gorm.DB, creditAccount *entities.CreditAccount, amount float64) error {
	numInstallments := s.calculateNumberOfInstallments(creditAccount)
	installmentAmount := s.calculateInstallmentAmount(creditAccount, amount, numInstallments)
	dueDate := s.calculateFirstInstallmentDueDate(creditAccount)

	for i := 0; i < numInstallments; i++ {
		installment := &entities.Installment{
			CreditAccountID: creditAccount.ID,
			DueDate:         dueDate.AddDate(0, i, 0), // Add a month to the due date
			Amount:          installmentAmount,
			Status:          enums.Pending,
		}
		if err := tx.Create(&installment).Error; err != nil {
			return fmt.Errorf("error creating installment: %w", err)
		}
	}

	return nil
}

// calculateNumberOfInstallments calculates the number of installments
func (s *purchaseService) calculateNumberOfInstallments(creditAccount *entities.CreditAccount) int {
	// Assuming the duration of long-term credit is always 12 months (1 year)
	creditDurationMonths := 12

	// Subtract the grace period (if any)
	numInstallments := creditDurationMonths - creditAccount.GracePeriod

	// Ensure the number of installments is at least 1
	if numInstallments < 1 {
		numInstallments = 1
	}

	return numInstallments
}

// calculateInstallmentAmount calculates the amount of each installment
func (s *purchaseService) calculateInstallmentAmount(creditAccount *entities.CreditAccount, amount float64, numInstallments int) float64 {
	// Calculate the periodic interest rate
	var periodicInterestRate float64
	if creditAccount.InterestType == enums.Nominal {
		// Nominal Rate: divide the annual rate by the number of periods in a year
		periodicInterestRate = creditAccount.InterestRate / 1200 // 12 months in a year
	} else if creditAccount.InterestType == enums.Effective {
		// Effective Rate: convert the effective annual rate to a periodic rate
		periodicInterestRate = math.Pow(1+(creditAccount.InterestRate/100), 1.0/12) - 1
	}

	// Calculate the installment amount using the annuity formula
	installmentAmount := amount * (periodicInterestRate * math.Pow(1+periodicInterestRate, float64(numInstallments))) / (math.Pow(1+periodicInterestRate, float64(numInstallments)) - 1)

	return installmentAmount
}

// calculateFirstInstallmentDueDate calculates the due date of the first installment
func (s *purchaseService) calculateFirstInstallmentDueDate(creditAccount *entities.CreditAccount) time.Time {
	today := time.Now()

	// Add the grace period to the current date
	firstDueDate := today.AddDate(0, creditAccount.GracePeriod, 0)

	// Adjust the due date to the day specified in MonthlyDueDate
	return time.Date(firstDueDate.Year(), firstDueDate.Month(), creditAccount.MonthlyDueDate, 0, 0, 0, 0, time.UTC)
}

// convertUserToClient converts a user to a client
func (s *purchaseService) convertUserToClient(tx *gorm.DB, userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("error retrieving user: %w", err)
	}

	if user.Rol == enums.USER {
		user.Rol = enums.CLIENT
		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("error updating user role: %w", err)
		}

		client := &entities.Client{
			UserID:   user.ID,
			IsActive: true,
		}
		if err := tx.Create(&client).Error; err != nil {
			return fmt.Errorf("error creating client record: %w", err)
		}
	}
	return nil
}

// GetClientBalance retrieves the current balance of a client's credit account
func (s *purchaseService) GetClientBalance(clientID uint) (float64, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByUserID(clientID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving credit account: %w", err)
	}
	return creditAccount.CurrentBalance, nil
}

// GetClientOverdueBalance retrieves the overdue balance of a client's credit account
func (s *purchaseService) GetClientOverdueBalance(clientID uint) (float64, error) {
	return s.creditAccountRepo.GetOverdueBalance(clientID)
}

// GetClientInstallments retrieves the installments of a client's credit account
func (s *purchaseService) GetClientInstallments(clientID uint) ([]response.InstallmentResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByUserID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	return s.installmentRepo.GetByCreditAccountID(creditAccount.ID)
}

// GetClientTransactions retrieves the transactions of a client's credit account
func (s *purchaseService) GetClientTransactions(clientID uint) ([]response.TransactionResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByUserID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	return s.transactionRepo.GetByCreditAccountID(creditAccount.ID)
}

// GetClientCreditAccount retrieves the credit account of a client
func (s *purchaseService) GetClientCreditAccount(clientID uint) (*entities.CreditAccount, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByUserID(clientID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	return creditAccount, nil
}
