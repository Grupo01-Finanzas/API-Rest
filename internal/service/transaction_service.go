package service

import (
	"ApiRestFinance/internal/model/entities/enums"
	"fmt"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/repository"
)

// TransactionService defines the interface for transaction service operations.
type TransactionService interface {
	CreateTransaction(req request.CreateTransactionRequest) (*response.TransactionResponse, error)
	GetTransactionByID(id uint) (*response.TransactionResponse, error)
	UpdateTransaction(id uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error)
	DeleteTransaction(id uint) error
	GetTransactionsByCreditAccountID(creditAccountID uint) ([]response.TransactionResponse, error)
}

type transactionService struct {
	transactionRepo   repository.TransactionRepository
	creditAccountRepo repository.CreditAccountRepository
}

// NewTransactionService creates a new instance of TransactionService.
func NewTransactionService(transactionRepo repository.TransactionRepository, creditAccountRepo repository.CreditAccountRepository) TransactionService {
	return &transactionService{
		transactionRepo:   transactionRepo,
		creditAccountRepo: creditAccountRepo,
	}
}

// CreateTransaction creates a new transaction, updating the credit account balance.
func (s *transactionService) CreateTransaction(req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	// Retrieve the credit account to update the balance
	creditAccount, err := s.creditAccountRepo.GetByID(req.CreditAccountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}

	// Create the transaction record
	transaction, err := s.transactionRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	// Update the credit account balance based on transaction type
	switch req.TransactionType {
	case enums.Purchase, enums.InterestAccrual, enums.LateFeeApplied: // Use enums.TransactionType
		creditAccount.CurrentBalance += req.Amount
	case enums.Payment, enums.EarlyPayment: // Use enums.TransactionType
		creditAccount.CurrentBalance -= req.Amount
	// Add other transaction types as needed
	default:
		return nil, fmt.Errorf("invalid transaction type: %s", req.TransactionType)
	}

	// Update the credit account in the database
	_, err = s.creditAccountRepo.Update(creditAccount.ID, request.UpdateCreditAccountRequest{
		CurrentBalance: creditAccount.CurrentBalance, // Assign to correct field
	})
	if err != nil {
		return nil, fmt.Errorf("error updating credit account balance: %w", err)
	}

	return transaction, nil
}

// GetTransactionByID retrieves a transaction by its ID.
func (s *transactionService) GetTransactionByID(id uint) (*response.TransactionResponse, error) {
	return s.transactionRepo.GetByID(id)
}

// UpdateTransaction updates an existing transaction and adjusts the credit account balance.
func (s *transactionService) UpdateTransaction(id uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error) {
	// Retrieve the original transaction
	originalTransaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transaction: %w", err)
	}

	// Retrieve the credit account to update the balance
	creditAccount, err := s.creditAccountRepo.GetByID(originalTransaction.CreditAccountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}

	// Adjust the credit account balance based on the original and new transaction amounts
	switch originalTransaction.TransactionType {
	case enums.Purchase, enums.InterestAccrual, enums.LateFeeApplied:
		creditAccount.CurrentBalance -= originalTransaction.Amount
	case enums.Payment, enums.EarlyPayment:
		creditAccount.CurrentBalance += originalTransaction.Amount
		// Add other transaction types as needed
	}

	if req.Amount > 0 { // Only adjust if a new amount is provided
		switch req.TransactionType {
		case enums.Purchase, enums.InterestAccrual, enums.LateFeeApplied:
			creditAccount.CurrentBalance += req.Amount
		case enums.Payment, enums.EarlyPayment:
			creditAccount.CurrentBalance -= req.Amount
			// Add other transaction types as needed
		}
	}

	// Update the transaction
	updatedTransaction, err := s.transactionRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("error updating transaction: %w", err)
	}

	// Update the credit account balance in the database
	_, err = s.creditAccountRepo.Update(creditAccount.ID, request.UpdateCreditAccountRequest{
		CurrentBalance: creditAccount.CurrentBalance,
	})
	if err != nil {
		return nil, fmt.Errorf("error updating credit account balance: %w", err)
	}

	return updatedTransaction, nil
}

// DeleteTransaction deletes a transaction and adjusts the credit account balance.
func (s *transactionService) DeleteTransaction(id uint) error {
	// Retrieve the transaction to be deleted
	transaction, err := s.transactionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("error retrieving transaction: %w", err)
	}

	// Retrieve the credit account to update the balance
	creditAccount, err := s.creditAccountRepo.GetByID(transaction.CreditAccountID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}

	// Adjust the credit account balance based on the deleted transaction
	switch transaction.TransactionType {
	case enums.Purchase, enums.InterestAccrual, enums.LateFeeApplied:
		creditAccount.CurrentBalance -= transaction.Amount
	case enums.Payment, enums.EarlyPayment:
		creditAccount.CurrentBalance += transaction.Amount
	// Add other transaction types as needed
	default:
		return fmt.Errorf("invalid transaction type: %s", transaction.TransactionType)
	}

	// Delete the transaction
	err = s.transactionRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}

	// Update the credit account balance in the database
	_, err = s.creditAccountRepo.Update(creditAccount.ID, request.UpdateCreditAccountRequest{
		CurrentBalance: creditAccount.CurrentBalance,
	})
	if err != nil {
		return fmt.Errorf("error updating credit account balance: %w", err)
	}

	return nil
}

// GetTransactionsByCreditAccountID retrieves transactions for a specific credit account.
func (s *transactionService) GetTransactionsByCreditAccountID(creditAccountID uint) ([]response.TransactionResponse, error) {
	return s.transactionRepo.GetByCreditAccountID(creditAccountID)
}
