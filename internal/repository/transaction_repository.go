package repository

import (
	"time"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"gorm.io/gorm"
)

// TransactionRepository defines the interface for transaction repository operations.
type TransactionRepository interface {
	Create(req request.CreateTransactionRequest) (*response.TransactionResponse, error)
	GetByID(id uint) (*response.TransactionResponse, error)
	Update(id uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error)
	Delete(id uint) error
	GetByCreditAccountID(creditAccountID uint) ([]response.TransactionResponse, error)
	GetTransactionsByDateRange(creditAccountID uint, startDate, endDate time.Time) ([]entities.Transaction, error)
	GetTransactionsByCreditAccountIDAndDateRange(creditAccountID uint, startDate, endDate time.Time) ([]entities.Transaction, error)
	GetBalanceBeforeDate(creditAccountID uint, beforeDate time.Time) (float64, error)
	GetTransactionsByCreditAccountID(creditAccountID uint) ([]entities.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new instance of transactionRepository.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create creates a new transaction.
func (r *transactionRepository) Create(req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	transaction := entities.Transaction{
		CreditAccountID: req.CreditAccountID,
		RecipientType:   req.RecipientType,
		RecipientID:     req.RecipientID,
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: time.Now(),
	}

	err := r.db.Create(&transaction).Error
	if err != nil {
		return nil, err
	}

	return getTransactionResponse(&transaction), nil
}

// GetByID retrieves a transaction by ID.
func (r *transactionRepository) GetByID(id uint) (*response.TransactionResponse, error) {
	var transaction entities.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}

	return getTransactionResponse(&transaction), nil
}

// Update updates an existing transaction.
func (r *transactionRepository) Update(id uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error) {
	var transaction entities.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}

	// Only update fields if they are provided in the request
	if req.Amount > 0 {
		transaction.Amount = req.Amount
	}
	if req.Description != "" {
		transaction.Description = req.Description
	}

	err = r.db.Save(&transaction).Error
	if err != nil {
		return nil, err
	}

	return getTransactionResponse(&transaction), nil
}

// Delete deletes a transaction.
func (r *transactionRepository) Delete(id uint) error {
	var transaction entities.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return err
	}

	return r.db.Delete(&transaction).Error
}

// GetByCreditAccountID retrieves all transactions for a credit account.
func (r *transactionRepository) GetByCreditAccountID(creditAccountID uint) ([]response.TransactionResponse, error) {
	var transactions []entities.Transaction
	err := r.db.Where("credit_account_id = ?", creditAccountID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	var transactionResponses []response.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, *getTransactionResponse(&transaction))
	}

	return transactionResponses, nil
}

// GetTransactionsByDateRange retrieves transactions for a credit account within a date range.
func (r *transactionRepository) GetTransactionsByDateRange(creditAccountID uint, startDate, endDate time.Time) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.Where("credit_account_id = ? AND transaction_date BETWEEN ? AND ?", creditAccountID, startDate, endDate).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func getTransactionResponse(transaction *entities.Transaction) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:              transaction.ID,
		CreditAccountID: transaction.CreditAccountID,
		TransactionType: transaction.TransactionType,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}
}

func (r *transactionRepository) GetTransactionsByCreditAccountIDAndDateRange(creditAccountID uint, startDate, endDate time.Time) ([]entities.Transaction, error) {
    var transactions []entities.Transaction
    db := r.db

    // Build the query
    query := db.Where("credit_account_id = ?", creditAccountID)

    // Add date range filters if provided
    if !startDate.IsZero() {
        query = query.Where("created_at >= ?", startDate)
    }
    if !endDate.IsZero() {
        query = query.Where("created_at <= ?", endDate)
    }

    err := query.Find(&transactions).Error
    if err != nil {
        return nil, err
    }

    return transactions, nil
}

func (r *transactionRepository) GetBalanceBeforeDate(creditAccountID uint, beforeDate time.Time) (float64, error) {
    var balance float64
    err := r.db.Model(&entities.Transaction{}).
        Select("SUM(CASE WHEN transaction_type = 'Payment' THEN amount ELSE -amount END) as balance"). // Adjust 'Payment' and transaction types as needed
        Where("credit_account_id = ? AND created_at < ?", creditAccountID, beforeDate).
        Scan(&balance).Error

    if err != nil {
        return 0, err
    }

    return balance, nil
}

// GetTransactionsByCreditAccountID retrieves all transactions for a given credit account ID.
func (r *transactionRepository) GetTransactionsByCreditAccountID(creditAccountID uint) ([]entities.Transaction, error) {
    var transactions []entities.Transaction
    err := r.db.Where("credit_account_id = ?", creditAccountID).Find(&transactions).Error
    if err != nil {
        return nil, err
    }
    return transactions, nil
}