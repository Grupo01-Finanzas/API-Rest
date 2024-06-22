package repository

import (
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreditAccountRepository defines operations for managing CreditAccount entities.
type CreditAccountRepository interface {
	CreateCreditAccount(creditAccount *entities.CreditAccount) error
	GetCreditAccountByID(creditAccountID uint) (*entities.CreditAccount, error)
	GetCreditAccountByClientID(clientID uint) (*entities.CreditAccount, error) 
	UpdateCreditAccount(creditAccount *entities.CreditAccount) error
	DeleteCreditAccount(creditAccountID uint) error
	GetCreditAccountsByEstablishmentID(establishmentID uint) ([]entities.CreditAccount, error) 
	ApplyInterest(creditAccount *entities.CreditAccount) error                    
	ApplyLateFee(creditAccount *entities.CreditAccount, daysOverdue int) error 
	GetOverdueCreditAccounts(establishmentID uint) ([]entities.CreditAccount, error)
	ProcessPurchase(creditAccount *entities.CreditAccount, amount float64, description string) error
	ProcessPayment(creditAccount *entities.CreditAccount, amount float64, description string) error
	CreateClientAndCreditAccount(user *entities.User, creditAccount *entities.CreditAccount) error
	DeleteClientAndCreditAccount(userID uint) error
	ProcessPurchaseTransaction(creditAccount *entities.CreditAccount, amount float64, description string) error
}

type creditAccountRepository struct {
	db *gorm.DB
	userRepo UserRepository
}

// NewCreditAccountRepository creates a new CreditAccountRepository instance.
func NewCreditAccountRepository(db *gorm.DB, userRepo UserRepository) CreditAccountRepository {
	return &creditAccountRepository{db: db, userRepo: userRepo}
}

// CreateCreditAccount creates a new credit account in the database.
func (r *creditAccountRepository) CreateCreditAccount(creditAccount *entities.CreditAccount) error {
	return r.db.Create(creditAccount).Error
}

// GetCreditAccountByID retrieves a credit account by its ID, including the Establishment.
func (r *creditAccountRepository) GetCreditAccountByID(creditAccountID uint) (*entities.CreditAccount, error) {
	var creditAccount entities.CreditAccount
	err := r.db.Preload("Client").Preload("Establishment").First(&creditAccount, creditAccountID).Error 
	if err != nil {
		return nil, err
	}
	return &creditAccount, nil
}

// GetCreditAccountByClientID retrieves a credit account by its client ID.
func (r *creditAccountRepository) GetCreditAccountByClientID(clientID uint) (*entities.CreditAccount, error) {
    var creditAccount entities.CreditAccount
    err := r.db.Where("client_id = ?", clientID).Preload("Client").Preload("Establishment").First(&creditAccount).Error 
    if err != nil {
        return nil, err
    }
    return &creditAccount, nil
}


// UpdateCreditAccount updates an existing credit account in the database.
func (r *creditAccountRepository) UpdateCreditAccount(creditAccount *entities.CreditAccount) error {
	return r.db.Save(creditAccount).Error
}

// DeleteCreditAccount deletes a credit account from the database.
func (r *creditAccountRepository) DeleteCreditAccount(creditAccountID uint) error {
	return r.db.Delete(&entities.CreditAccount{}, creditAccountID).Error
}

// GetCreditAccountsByEstablishmentID retrieves all credit accounts for an establishment.
func (r *creditAccountRepository) GetCreditAccountsByEstablishmentID(establishmentID uint) ([]entities.CreditAccount, error) {
    var creditAccounts []entities.CreditAccount
    err := r.db.Preload("Client").Preload("Establishment").Where("establishment_id = ?", establishmentID).Find(&creditAccounts).Error 
    if err != nil {
        return nil, err
    }
    return creditAccounts, nil
}


// ApplyInterest calculates and applies interest to a credit account.
func (r *creditAccountRepository) ApplyInterest(creditAccount *entities.CreditAccount) error {
	if creditAccount.CurrentBalance == 0 || 
	   time.Now().Before(creditAccount.LastInterestAccrualDate.AddDate(0, 1, 0)) {
		return nil 
	}

	interest := calculateInterest(*creditAccount)
	creditAccount.CurrentBalance += interest
	creditAccount.LastInterestAccrualDate = time.Now()

	return r.db.Save(creditAccount).Error
}

// ApplyLateFee applies late fee to a credit account.
func (r *creditAccountRepository) ApplyLateFee(creditAccount *entities.CreditAccount, daysOverdue int) error {
	if daysOverdue <= 0 {
		return nil
	}

	lateFee := creditAccount.CurrentBalance * (creditAccount.Establishment.LateFeePercentage / 100) 
	creditAccount.CurrentBalance += lateFee

	return r.db.Save(creditAccount).Error
}

// GetOverdueCreditAccounts gets all overdue credit accounts for an establishment.
func (r *creditAccountRepository) GetOverdueCreditAccounts(establishmentID uint) ([]entities.CreditAccount, error) {
	today := time.Now()
	var overdueAccounts []entities.CreditAccount
	err := r.db.Preload("Client").Preload("Establishment").Where("establishment_id = ? AND monthly_due_date < ? AND current_balance > 0", establishmentID, today.Day()).Find(&overdueAccounts).Error 
	if err != nil {
		return nil, err
	}
	return overdueAccounts, nil
}

func (r *creditAccountRepository) ProcessPurchase(creditAccount *entities.CreditAccount, amount float64, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if creditAccount.IsBlocked {
			return errors.New("credit account is blocked, cannot process purchase")
		}

		if creditAccount.CurrentBalance+amount > creditAccount.CreditLimit {
			return errors.New("purchase exceeds credit limit")
		}

		// Create the purchase transaction
		transaction := entities.Transaction{
			CreditAccountID: creditAccount.ID,
			TransactionType: enums.Purchase,
			Amount:          amount,
			Description:     description,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("error creating purchase transaction: %w", err)
		}

		// Update the credit account balance
		creditAccount.CurrentBalance += amount
		if err := tx.Save(creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		return nil
	})
}

func (r *creditAccountRepository) ProcessPayment(creditAccount *entities.CreditAccount, amount float64, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Retrieve the credit account for update, locking the row
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(creditAccount, creditAccount.ID).Error; err != nil {
			return fmt.Errorf("error retrieving credit account for payment: %w", err)
		}

		if amount > creditAccount.CurrentBalance {
			return fmt.Errorf("payment amount exceeds current balance: %.2f", creditAccount.CurrentBalance) 
		}

		transaction := entities.Transaction{
			CreditAccountID: creditAccount.ID,
			TransactionType: enums.Payment,
			Amount:          amount,
			Description:     description,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("error creating payment transaction: %w", err)
		}

		creditAccount.CurrentBalance -= amount
		if err := tx.Save(creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		if creditAccount.IsBlocked && creditAccount.CurrentBalance <= 0 {
			creditAccount.IsBlocked = false
			if err := tx.Save(creditAccount).Error; err != nil {
				return fmt.Errorf("error unblocking credit account: %w", err)
			}
		}

		return nil 
	})
}

// calculateInterest calculates the interest for a credit account based on its type and interest type
func calculateInterest(creditAccount entities.CreditAccount) float64 {
	var interest float64
	principal := creditAccount.CurrentBalance
	annualRate := creditAccount.InterestRate / 100 
	timeInYears := 1.0 / 12.0 

	if creditAccount.InterestType == enums.Nominal {
		interest = principal * annualRate * timeInYears
	} else if creditAccount.InterestType == enums.Effective {
		interest = principal * (math.Pow(1+annualRate, timeInYears) - 1)
	}
	return interest
}

func (r *creditAccountRepository) DeleteCreditAccountInTransaction(tx *gorm.DB, creditAccountID uint) error {
    return tx.Delete(&entities.CreditAccount{}, creditAccountID).Error
}

// CreateClientAndCreditAccount creates a new client user and their credit account in a transaction.
func (r *creditAccountRepository) CreateClientAndCreditAccount(user *entities.User, creditAccount *entities.CreditAccount) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		creditAccount.ClientID = user.ID // Associate the user ID after creation
		if err := tx.Create(creditAccount).Error; err != nil {
			return fmt.Errorf("error creating credit account: %w", err)
		}

		return nil // Transaction successful
	})
}

func (r *creditAccountRepository) DeleteClientAndCreditAccount(userID uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. Get the CreditAccount ID 
        creditAccount, err := r.GetCreditAccountByClientID(userID)
        if err != nil {
            return fmt.Errorf("error retrieving credit account: %w", err)
        }

        // 2. Delete the Credit Account
        if err := r.DeleteCreditAccountInTransaction(tx, creditAccount.ID); err != nil {
            return fmt.Errorf("error deleting credit account: %w", err)
        }

        // 3. Delete the User
        // You can access the userRepo from here if you pass it during initialization
        // For example, if your creditAccountRepository has a userRepo field:
        if err := r.userRepo.DeleteUser(userID); err != nil {
            return fmt.Errorf("error deleting user: %w", err)
        }

        return nil // Transaction successful
    })
}

// ProcessPurchaseTransaction handles the purchase logic within a transaction.
func (r *creditAccountRepository) ProcessPurchaseTransaction(creditAccount *entities.CreditAccount, amount float64, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if creditAccount.IsBlocked {
			return errors.New("credit account is blocked, cannot process purchase")
		}

		if creditAccount.CurrentBalance+amount > creditAccount.CreditLimit {
			return errors.New("purchase exceeds credit limit")
		}

		// Create the purchase transaction
		transaction := entities.Transaction{
			CreditAccountID: creditAccount.ID,
			TransactionType: enums.Purchase,
			Amount:          amount,
			Description:     description,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("error creating purchase transaction: %w", err)
		}

		// Update the credit account's current balance
		creditAccount.CurrentBalance += amount
		if err := tx.Save(creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		return nil
	})
}