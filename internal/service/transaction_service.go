package service

import (
	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"ApiRestFinance/internal/repository"
	"ApiRestFinance/internal/util"
	"errors"
	"fmt"
	"time"
)

// TransactionService handles transaction-related operations.
type TransactionService interface {
	CreateTransaction(req request.CreateTransactionRequest) (*response.TransactionResponse, error)
	GetTransactionByID(id uint) (*response.TransactionResponse, error)
	GetTransactionsByCreditAccountID(creditAccountID uint) ([]response.TransactionResponse, error)
	UpdateTransaction(id uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error)
	DeleteTransaction(id uint) error
	ConfirmPayment(transactionID uint, confirmationCode string) error
}

type transactionService struct {
	transactionRepo   repository.TransactionRepository
	creditAccountRepo repository.CreditAccountRepository
}

// NewTransactionService creates a new TransactionService instance.
func NewTransactionService(transactionRepo repository.TransactionRepository, creditAccountRepo repository.CreditAccountRepository) TransactionService {
	return &transactionService{
		transactionRepo:   transactionRepo,
		creditAccountRepo: creditAccountRepo,
	}
}

func (s *transactionService) CreateTransaction(req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(req.CreditAccountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return nil, errors.New("credit account not found")
	}

	if req.Amount <= 0 {
		return nil, errors.New("transaction amount must be greater than zero")
	}

	var paymentCode string
	if req.PaymentMethod != enums.CASH {
		paymentCode = util.GeneratePaymentCode()
	}

	transaction := entities.Transaction{
		CreditAccountID: creditAccount.ID,
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: time.Now(),
		PaymentMethod:   req.PaymentMethod,
		PaymentCode:     paymentCode,
		PaymentStatus:   enums.PENDING,
	}

	if err := s.transactionRepo.CreateTransaction(&transaction, creditAccount); err != nil {
		return nil, fmt.Errorf("error processing transaction: %w", err)
	}
	return transactionToResponse(&transaction), nil
}

func (s *transactionService) ConfirmPayment(transactionID uint, confirmationCode string) error {
	transaction, err := s.transactionRepo.GetTransactionByID(transactionID)
	if err != nil {
		return fmt.Errorf("error retrieving transaction: %w", err)
	}
	if transaction == nil {
		return errors.New("transaction not found")
	}

	// Check if the transaction is pending and the payment method is not cash
	if transaction.PaymentStatus != enums.PENDING || transaction.PaymentMethod == enums.CASH {
		return errors.New("transaction cannot be confirmed")
	}

	// Validate the confirmation code against the generated PaymentCode
	if transaction.PaymentCode != confirmationCode {
		transaction.PaymentStatus = enums.FAILED
		if err := s.transactionRepo.UpdateTransaction(transaction, nil); err != nil {
			return fmt.Errorf("error updating transaction: %w", err)
		}

		if transaction.PaymentCode == "" {
			transaction.PaymentCode = util.GeneratePaymentCode()
		}
		return errors.New("invalid confirmation code")
	}

	// Update the transaction status to SUCCESS
	transaction.PaymentStatus = enums.SUCCESS
	transaction.ConfirmationCode = confirmationCode

	return s.transactionRepo.UpdateTransaction(transaction, nil)
}

func (s *transactionService) GetTransactionByID(id uint) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transaction: %w", err)
	}
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	return transactionToResponse(transaction), nil
}

func (s *transactionService) GetTransactionsByCreditAccountID(creditAccountID uint) ([]response.TransactionResponse, error) {
	transactions, err := s.transactionRepo.GetTransactionsByCreditAccountID(creditAccountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions: %w", err)
	}

	var transactionResponses []response.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, *transactionToResponse(&transaction))
	}

	return transactionResponses, nil
}

func (s *transactionService) UpdateTransaction(id uint, req request.UpdateTransactionRequest) (*response.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transaction: %w", err)
	}
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	// Retrieve the credit account
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(transaction.CreditAccountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return nil, errors.New("credit account not found")
	}

	// Update transaction details
	if req.Amount > 0 {
		transaction.Amount = req.Amount
	}
	if req.Description != "" {
		transaction.Description = req.Description
	}
	if req.TransactionType != "" {
		transaction.TransactionType = req.TransactionType
	}

	// Update the transaction and credit account balance in a transaction
	if err := s.transactionRepo.UpdateTransaction(transaction, creditAccount); err != nil {
		return nil, fmt.Errorf("error updating transaction: %w", err)
	}

	return transactionToResponse(transaction), nil
}

func (s *transactionService) DeleteTransaction(id uint) error {
	transaction, err := s.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return fmt.Errorf("error retrieving transaction: %w", err)
	}
	if transaction == nil {
		return errors.New("transaction not found")
	}

	// Retrieve the credit account to adjust the balance
	creditAccount, err := s.creditAccountRepo.GetCreditAccountByID(transaction.CreditAccountID)
	if err != nil {
		return fmt.Errorf("error retrieving credit account: %w", err)
	}
	if creditAccount == nil {
		return errors.New("credit account not found")
	}

	// Delete the transaction and update the credit account balance
	if err := s.transactionRepo.DeleteTransaction(id, creditAccount); err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}

	return nil
}

func transactionToResponse(transaction *entities.Transaction) *response.TransactionResponse {
	return &response.TransactionResponse{
		ID:              transaction.ID,
		CreditAccountID: transaction.CreditAccountID,
		TransactionType: transaction.TransactionType,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate,
		PaymentMethod:   transaction.PaymentMethod,
		PaymentCode:     transaction.PaymentCode,
		PaymentStatus:   transaction.PaymentStatus,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}
}
