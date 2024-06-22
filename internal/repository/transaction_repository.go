package repository

import (
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TransactionRepository defines operations for managing Transaction entities.
type TransactionRepository interface {
	CreateTransaction(transaction *entities.Transaction, creditAccount *entities.CreditAccount) error
	GetTransactionByID(transactionID uint) (*entities.Transaction, error)
	GetTransactionsByCreditAccountID(creditAccountID uint) ([]entities.Transaction, error)
	UpdateTransaction(transaction *entities.Transaction, creditAccount *entities.CreditAccount) error
	DeleteTransaction(transactionID uint, creditAccount *entities.CreditAccount) error
	CreateTransactionInTx(tx *gorm.DB, transaction *entities.Transaction) error
	UpdateTransactionInTx(tx *gorm.DB, transaction *entities.Transaction) error
	DeleteTransactionInTx(tx *gorm.DB, transactionID uint) error
	GetTransactionsByCreditAccountIDAndDateRange(creditAccountID uint, startDate, endDate time.Time) ([]entities.Transaction, error)
	GetBalanceBeforeDate(creditAccountID uint, beforeDate time.Time) (float64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new TransactionRepository instance.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// CreateTransaction creates a new transaction and updates the credit account balance in a transaction.
func (r *transactionRepository) CreateTransaction(transaction *entities.Transaction, creditAccount *entities.CreditAccount) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transaction).Error; err != nil {
			return fmt.Errorf("error creating transaction: %w", err)
		}

		// Update the credit account balance based on the transaction type
		switch transaction.TransactionType {
		case enums.Purchase:
			creditAccount.CurrentBalance += transaction.Amount
		case enums.Payment:
			if transaction.Amount > creditAccount.CurrentBalance {
				return fmt.Errorf("payment amount exceeds current balance: %.2f", creditAccount.CurrentBalance)
			}
			creditAccount.CurrentBalance -= transaction.Amount

			// Unblock the account if it was blocked and the balance is zero or less
			if creditAccount.IsBlocked && creditAccount.CurrentBalance <= 0 {
				creditAccount.IsBlocked = false
			}
		default:
			return errors.New("invalid transaction type")
		}

		// Save the updated credit account
		if err := tx.Save(creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		return nil
	})
}

// UpdateTransaction updates a transaction and adjusts the credit account balance in a transaction.
func (r *transactionRepository) UpdateTransaction(transaction *entities.Transaction, creditAccount *entities.CreditAccount) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Reverse the effect of the original transaction
		switch transaction.TransactionType {
		case enums.Purchase:
			creditAccount.CurrentBalance -= transaction.Amount
		case enums.Payment:
			creditAccount.CurrentBalance += transaction.Amount
		default:
			return errors.New("invalid transaction type")
		}

		// Update transaction details (no changes here)
		// ...

		// Apply the effect of the updated transaction
		switch transaction.TransactionType {
		case enums.Purchase:
			creditAccount.CurrentBalance += transaction.Amount
		case enums.Payment:
			if transaction.Amount > creditAccount.CurrentBalance {
				return fmt.Errorf("payment amount exceeds current balance: %.2f", creditAccount.CurrentBalance)
			}
			creditAccount.CurrentBalance -= transaction.Amount

			if creditAccount.IsBlocked && creditAccount.CurrentBalance <= 0 {
				creditAccount.IsBlocked = false
			}
		default:
			return errors.New("invalid transaction type")
		}

		// Save the updated transaction and credit account
		if err := tx.Save(transaction).Error; err != nil {
			return fmt.Errorf("error updating transaction: %w", err)
		}
		if err := tx.Save(creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		return nil
	})
}

// DeleteTransaction deletes a transaction and adjusts the credit account balance in a transaction.
func (r *transactionRepository) DeleteTransaction(transactionID uint, creditAccount *entities.CreditAccount) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Retrieve the transaction for deletion
		var transaction entities.Transaction
		if err := tx.First(&transaction, transactionID).Error; err != nil {
			return fmt.Errorf("error retrieving transaction: %w", err)
		}

		// Reverse the effect of the transaction on the credit account balance
		switch transaction.TransactionType {
		case enums.Purchase:
			creditAccount.CurrentBalance -= transaction.Amount
		case enums.Payment:
			creditAccount.CurrentBalance += transaction.Amount
		default:
			return errors.New("invalid transaction type")
		}

		// Delete the transaction
		if err := tx.Delete(&transaction).Error; err != nil {
			return fmt.Errorf("error deleting transaction: %w", err)
		}

		// Save the updated credit account balance
		if err := tx.Save(creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		return nil
	})
}

// GetTransactionByID retrieves a transaction by its ID.
func (r *transactionRepository) GetTransactionByID(transactionID uint) (*entities.Transaction, error) {
	var transaction entities.Transaction
	err := r.db.First(&transaction, transactionID).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByCreditAccountID retrieves all transactions for a specific credit account.
func (r *transactionRepository) GetTransactionsByCreditAccountID(creditAccountID uint) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.Where("credit_account_id = ?", creditAccountID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) CreateTransactionInTx(tx *gorm.DB, transaction *entities.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepository) UpdateTransactionInTx(tx *gorm.DB, transaction *entities.Transaction) error {
	return tx.Save(transaction).Error
}

func (r *transactionRepository) DeleteTransactionInTx(tx *gorm.DB, transactionID uint) error {
	return tx.Delete(&entities.Transaction{}, transactionID).Error
}

// GetTransactionsByCreditAccountIDAndDateRange retrieves transactions for a credit account within a given date range.
func (r *transactionRepository) GetTransactionsByCreditAccountIDAndDateRange(creditAccountID uint, startDate, endDate time.Time) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	db := r.db.Where("credit_account_id = ?", creditAccountID)

	if !startDate.IsZero() {
		db = db.Where("transaction_date >= ?", startDate)
	}
	if !endDate.IsZero() {
		db = db.Where("transaction_date <= ?", endDate)
	}

	err := db.Find(&transactions).Error
	return transactions, err
}

// GetBalanceBeforeDate retrieves the balance of a credit account before a specified date.
func (r *transactionRepository) GetBalanceBeforeDate(creditAccountID uint, beforeDate time.Time) (float64, error) {
	var balance float64
	err := r.db.Model(&entities.Transaction{}).
		Select("SUM(CASE WHEN transaction_type = ? THEN -amount ELSE amount END) as balance", enums.Payment).
		Where("credit_account_id = ? AND transaction_date < ?", creditAccountID, beforeDate).
		Scan(&balance).Error

	if err != nil {
		return 0, fmt.Errorf("error getting balance before date: %w", err)
	}

	return balance, nil
}
