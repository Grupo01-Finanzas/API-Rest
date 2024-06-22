package service

import "errors"

// Define custom errors
var (
	ErrCreditAccountNotFound  = errors.New("credit account not found")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInsufficientBalance    = errors.New("insufficient balance")
	ErrInvalidFileType        = errors.New("invalid file type. Only images are allowed")
	ErrFileSizeTooLarge       = errors.New("file size too large")
)
