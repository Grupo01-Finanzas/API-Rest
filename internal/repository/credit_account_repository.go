package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"math"
	"time"

	"ApiRestFinance/internal/model/dto/request"
	"ApiRestFinance/internal/model/dto/response"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
)

// CreditAccountRepository defines the interface for credit account repository operations.
type CreditAccountRepository interface {
	Create(creditAccount *entities.CreditAccount) error
	GetByID(id uint) (*response.CreditAccountResponse, error)
	Update(id uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error)
	Delete(id uint) error
	GetByEstablishmentID(establishmentID uint) ([]response.CreditAccountResponse, error)
	GetByClientID(clientID uint) ([]response.CreditAccountResponse, error)
	ApplyInterest(creditAccountID uint) error
	ApplyLateFee(creditAccountID uint) error
	GetOverdueAccounts(establishmentID uint) ([]response.CreditAccountResponse, error)
	ExistsByClientAndEstablishment(clientID uint, establishmentID uint) (bool, error)
	CreateCreditRequest(creditRequest *entities.CreditRequest) error
	GetCreditRequestByID(id uint) (*entities.CreditRequest, error)
	UpdateCreditRequest(creditRequest *entities.CreditRequest) error
	ProcessPurchase(creditAccountID uint, amount float64, description string) error
	ProcessPayment(creditAccountID uint, amount float64, description string) error
	ApproveCreditRequest(creditRequest *entities.CreditRequest) (*response.CreditAccountResponse, error)
	GetPendingCreditRequests(establishmentID uint) ([]entities.CreditRequest, error)
	AssignCreditAccountToClient(creditAccountID, clientID uint) error
	GetOverdueBalance(userID uint) (float64, error)
	CalculateInterest(creditAccount entities.CreditAccount) float64
	GetCreditAccountByUserID(userID uint) (*entities.CreditAccount, error)
}

type creditAccountRepository struct {
	db              *gorm.DB
	installmentRepo InstallmentRepository
}

// NewCreditAccountRepository creates a new instance of creditAccountRepository.
func NewCreditAccountRepository(db *gorm.DB) CreditAccountRepository {
	return &creditAccountRepository{db: db}
}

// Create creates a new credit account.
func (r *creditAccountRepository) Create(creditAccount *entities.CreditAccount) error {
	return r.db.Create(creditAccount).Error
}

// GetByID retrieves a credit account by ID.
func (r *creditAccountRepository) GetByID(id uint) (*response.CreditAccountResponse, error) {
	var creditAccount entities.CreditAccount
	err := r.db.Preload("Establishment").Preload("Client").First(&creditAccount, id).Error
	if err != nil {
		return nil, err
	}

	return getCreditAccountResponse(&creditAccount), nil
}

// Update updates an existing credit account.
func (r *creditAccountRepository) Update(id uint, req request.UpdateCreditAccountRequest) (*response.CreditAccountResponse, error) {
	var creditAccount entities.CreditAccount
	err := r.db.First(&creditAccount, id).Error
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

	err = r.db.Save(&creditAccount).Error
	if err != nil {
		return nil, err
	}

	return getCreditAccountResponse(&creditAccount), nil
}

// Delete deletes a credit account.
func (r *creditAccountRepository) Delete(id uint) error {
	var creditAccount entities.CreditAccount
	err := r.db.First(&creditAccount, id).Error
	if err != nil {
		return err
	}

	return r.db.Delete(&creditAccount).Error
}

// GetByEstablishmentID retrieves all credit accounts for an establishment.
func (r *creditAccountRepository) GetByEstablishmentID(establishmentID uint) ([]response.CreditAccountResponse, error) {
	var creditAccounts []entities.CreditAccount
	err := r.db.Where("establishment_id = ?", establishmentID).Preload("Client").Find(&creditAccounts).Error
	if err != nil {
		return nil, err
	}

	var creditAccountResponses []response.CreditAccountResponse
	for _, creditAccount := range creditAccounts {
		creditAccountResponses = append(creditAccountResponses, *getCreditAccountResponse(&creditAccount))
	}

	return creditAccountResponses, nil
}

// GetByClientID retrieves all credit accounts for a client.
func (r *creditAccountRepository) GetByClientID(clientID uint) ([]response.CreditAccountResponse, error) {
	var creditAccounts []entities.CreditAccount
	err := r.db.Where("client_id = ?", clientID).Preload("Establishment").Find(&creditAccounts).Error
	if err != nil {
		return nil, err
	}

	var creditAccountResponses []response.CreditAccountResponse
	for _, creditAccount := range creditAccounts {
		creditAccountResponses = append(creditAccountResponses, *getCreditAccountResponse(&creditAccount))
	}

	return creditAccountResponses, nil
}

// ApplyInterest calculates and applies interest to a credit account.
func (r *creditAccountRepository) ApplyInterest(creditAccountID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Retrieve the credit account
		var creditAccount entities.CreditAccount
		if err := tx.Preload("Establishment").First(&creditAccount, creditAccountID).Error; err != nil {
			return fmt.Errorf("error retrieving credit account: %w", err)
		}

		// 2. Check if interest needs to be applied
		if creditAccount.CurrentBalance == 0 || time.Now().Before(creditAccount.LastInterestAccrualDate.AddDate(0, 1, 0)) {
			return nil // No balance or interest already applied this month
		}

		// 3. Calculate interest based on the interest type (Nominal or Effective)
		interest := r.CalculateInterest(creditAccount)

		// 4. Create a transaction for the interest
		interestTransaction := entities.Transaction{
			CreditAccountID: creditAccountID,
			RecipientType:   enums.RolClient,
			RecipientID:     creditAccount.ClientID,
			TransactionType: enums.InterestAccrual,
			Amount:          interest,
			Description:     "Monthly Interest",
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&interestTransaction).Error; err != nil {
			return fmt.Errorf("error creating interest transaction: %w", err)
		}

		// 5. Update the credit account balance and last interest accrual date
		creditAccount.CurrentBalance += interest
		creditAccount.LastInterestAccrualDate = time.Now()
		if err := tx.Save(&creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		// 6. Create a credit account history record
		historyEntry := entities.CreditAccountHistory{
			CreditAccountID: creditAccountID,
			TransactionDate: time.Now(),
			TransactionType: enums.InterestAccrual,
			Amount:          interest,
			Balance:         creditAccount.CurrentBalance,
			Description:     "Monthly Interest Accrued",
		}
		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creating credit account history record: %w", err)
		}

		return nil
	})
}

// ApplyLateFee applies a late fee to a credit account based on the configured rules.
func (r *creditAccountRepository) ApplyLateFee(creditAccountID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Retrieve the credit account with LateFeeRule and Establishment
		var creditAccount entities.CreditAccount
		if err := tx.Preload("LateFeeRule").Preload("Establishment.LateFeeRules").First(&creditAccount, creditAccountID).Error; err != nil {
			return fmt.Errorf("error retrieving credit account: %w", err)
		}

		// 2. Calculate the number of days overdue
		daysOverdue := calculateDaysOverdue(creditAccount.MonthlyDueDate)
		if daysOverdue <= 0 {
			return nil // Account is not overdue
		}

		// 3. Find the applicable late fee rule
		var lateFeeRule *entities.LateFeeRule
		if creditAccount.LateFeeRule != nil {
			lateFeeRule = creditAccount.LateFeeRule
		} else {
			// Find a rule from the Establishment's rules based on days overdue
			for _, rule := range creditAccount.Establishment.LateFeeRules {
				if daysOverdue >= rule.DaysOverdueMin && daysOverdue <= rule.DaysOverdueMax {
					lateFeeRule = &rule
					break
				}
			}
		}

		// 4. Handle case where no applicable LateFeeRule is found
		if lateFeeRule == nil {
			// You might have a default rule or skip applying a late fee
			return errors.New("no late fee rule found for this credit account")
		}

		// 5. Calculate the late fee amount
		lateFeeAmount := calculateLateFee(creditAccount, *lateFeeRule, daysOverdue)

		// 6. Create a late fee record
		lateFee := entities.LateFee{
			CreditAccountID: creditAccountID,
			Amount:          lateFeeAmount,
			AppliedDate:     time.Now(),
		}
		if err := tx.Create(&lateFee).Error; err != nil {
			return fmt.Errorf("error creating late fee record: %w", err)
		}

		// 7. Update the credit account balance
		creditAccount.CurrentBalance += lateFeeAmount
		if err := tx.Save(&creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		// 8. Create a credit account history record
		historyEntry := entities.CreditAccountHistory{
			CreditAccountID: creditAccountID,
			TransactionDate: time.Now(),
			TransactionType: enums.LateFeeApplied,
			Amount:          lateFeeAmount,
			Balance:         creditAccount.CurrentBalance,
			Description:     "Late Payment Fee Applied",
		}
		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creating credit account history record: %w", err)
		}

		return nil
	})
}

// GetOverdueAccounts retrieves all overdue credit accounts for an establishment.
func (r *creditAccountRepository) GetOverdueAccounts(establishmentID uint) ([]response.CreditAccountResponse, error) {
	var overdueAccounts []entities.CreditAccount
	today := time.Now()
	err := r.db.Where("establishment_id = ? AND monthly_due_date < ? AND current_balance > 0", establishmentID, today.Day()).
		Preload("Client").
		Find(&overdueAccounts).Error
	if err != nil {
		return nil, err
	}

	var overdueAccountResponses []response.CreditAccountResponse
	for _, account := range overdueAccounts {
		overdueAccountResponses = append(overdueAccountResponses, *getCreditAccountResponse(&account))
	}

	return overdueAccountResponses, nil
}

// Helper functions for calculations

func (r *creditAccountRepository) CalculateInterest(creditAccount entities.CreditAccount) float64 {
	if creditAccount.CurrentBalance == 0 {
		return 0 // There is no interest to calculate
	}

	today := time.Now()

	// Calculate the number of days since the last interest accrual
	daysSinceLastAccrual := today.Sub(creditAccount.LastInterestAccrualDate).Hours() / 24

	if creditAccount.CreditType == enums.ShortTerm {
		// Short-Term Interest Calculation
		if creditAccount.InterestType == enums.Nominal {
			// Nominal Interest
			// Fórmula: Interest = Balance * (interest rate/100) * (Time in year)
			return creditAccount.CurrentBalance * (creditAccount.InterestRate / 100) * (daysSinceLastAccrual / 365)
		} else if creditAccount.InterestType == enums.Effective {
			// Effective Interest
			// Fórmula: Interest = Balance * ((1 + interest rate/100)^(Time in year) - 1)
			return creditAccount.CurrentBalance * (math.Pow(1+(creditAccount.InterestRate/100), daysSinceLastAccrual/365) - 1)
		}
	} else if creditAccount.CreditType == enums.LongTerm {
		// Short-Term Interest Calculation (Dues)

		// Get all overdue installments
		installments, err := r.installmentRepo.GetOverdueInstallments(creditAccount.ID)
		if err != nil {
			// Handle error
			fmt.Printf("Error al obtener las cuotas: %v", err)
			return 0
		}

		totalInterest := 0.0
		for _, installment := range installments {
			// Calcular la cantidad de días de atraso para cada cuota
			daysOverdue := int(today.Sub(installment.DueDate).Hours() / 24)

			if daysOverdue > 0 { // Solo calcular interés en cuotas vencidas
				if creditAccount.InterestType == enums.Nominal {
					// Interés nominal
					dailyRate := creditAccount.InterestRate / 36500 // Tasa diaria nominal
					totalInterest += installment.Amount * dailyRate * float64(daysOverdue)
				} else if creditAccount.InterestType == enums.Effective {
					// Interés efectivo
					dailyRate := math.Pow(1+(creditAccount.InterestRate/100), 1.0/365) - 1 // Tasa diaria efectiva
					totalInterest += installment.Amount*math.Pow(1+dailyRate, float64(daysOverdue)) - installment.Amount
				}
			}
		}

		return totalInterest
	}

	return 0 // Si el tipo de crédito no es válido, no se calcula el interés
}

func calculateDaysOverdue(monthlyDueDate int) int {
	today := time.Now()
	dueDate := time.Date(today.Year(), today.Month(), monthlyDueDate, 0, 0, 0, 0, time.UTC)

	if today.Before(dueDate) {
		return 0 // Not overdue
	}

	return int(today.Sub(dueDate).Hours() / 24)
}

func calculateLateFee(creditAccount entities.CreditAccount, rule entities.LateFeeRule, daysOverdue int) float64 {
	if daysOverdue < rule.DaysOverdueMin || daysOverdue > rule.DaysOverdueMax {
		return 0 // Rule does not apply
	}

	if rule.FeeType == enums.Percentage {
		return rule.FeeValue / 100 * creditAccount.CurrentBalance // Percentage of balance
	} else { // enums.FixedAmount
		return rule.FeeValue // Fixed amount
	}
}

func getCreditAccountResponse(creditAccount *entities.CreditAccount) *response.CreditAccountResponse {
	return &response.CreditAccountResponse{
		ID:              creditAccount.ID,
		EstablishmentID: creditAccount.EstablishmentID,
		ClientID:        creditAccount.ClientID,
		CreditLimit:     creditAccount.CreditLimit,
		CurrentBalance:  creditAccount.CurrentBalance,
		MonthlyDueDate:  creditAccount.MonthlyDueDate,
		InterestRate:    creditAccount.InterestRate,
		InterestType:    creditAccount.InterestType,
		CreditType:      creditAccount.CreditType,
		GracePeriod:     creditAccount.GracePeriod,
		IsBlocked:       creditAccount.IsBlocked,
		CreatedAt:       creditAccount.CreatedAt,
		UpdatedAt:       creditAccount.UpdatedAt,
		Client:          creditAccount.Client,
	}
}

/*func getClientResponse(client *entities.Client) *response.ClientResponse {
	if client == nil {
		return nil
	}

	return &response.ClientResponse{
		ID:        client.ID,
		User:      client.User,
		IsActive:  client.IsActive,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}*/

func (r *creditAccountRepository) ExistsByClientAndEstablishment(clientID uint, establishmentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entities.CreditAccount{}).
		Where("client_id = ? AND establishment_id = ?", clientID, establishmentID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error checking for existing credit account: %w", err)
	}
	return count > 0, nil
}

func (r *creditAccountRepository) CreateCreditRequest(creditRequest *entities.CreditRequest) error {
	return r.db.Create(creditRequest).Error
}

// UpdateCreditRequest updates a credit request in the database.
func (r *creditAccountRepository) UpdateCreditRequest(creditRequest *entities.CreditRequest) error {
	return r.db.Save(creditRequest).Error
}

func (r *creditAccountRepository) ProcessPurchase(creditAccountID uint, amount float64, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Retrieve the credit account
		var creditAccount entities.CreditAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&creditAccount, creditAccountID).Error; err != nil {
			return fmt.Errorf("error retrieving credit account for purchase: %w", err)
		}

		// 2. Check if the account is blocked
		if creditAccount.IsBlocked {
			return errors.New("credit account is blocked, cannot process purchase")
		}

		// 3. Check if the purchase exceeds the credit limit
		if creditAccount.CurrentBalance+amount > creditAccount.CreditLimit {
			return errors.New("purchase exceeds credit limit")
		}

		// 4. Create the purchase transaction
		transaction := entities.Transaction{
			CreditAccountID: creditAccountID,
			RecipientType:   enums.RolClient,
			RecipientID:     creditAccount.ClientID,
			TransactionType: enums.Purchase,
			Amount:          amount,
			Description:     description,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("error creating purchase transaction: %w", err)
		}

		// 5. Update the credit account balance
		creditAccount.CurrentBalance += amount
		if err := tx.Save(&creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		// 6. Create a CreditAccountHistory record
		historyEntry := entities.CreditAccountHistory{
			CreditAccountID: creditAccountID,
			TransactionDate: time.Now(),
			TransactionType: enums.Purchase,
			Amount:          amount,
			Balance:         creditAccount.CurrentBalance,
			Description:     description,
		}
		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creating credit account history: %w", err)
		}

		return nil
	})
}

func (r *creditAccountRepository) ProcessPayment(creditAccountID uint, amount float64, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Retrieve the credit account
		var creditAccount entities.CreditAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&creditAccount, creditAccountID).Error; err != nil {
			return fmt.Errorf("error retrieving credit account for payment: %w", err)
		}

		// 2. Check if payment exceeds current balance
		if amount > creditAccount.CurrentBalance {
			return errors.New("payment amount exceeds current balance")
		}

		// 3. Create the payment transaction
		transaction := entities.Transaction{
			CreditAccountID: creditAccountID,
			RecipientType:   enums.RolEstablishment, // Payment is to the establishment
			RecipientID:     creditAccount.EstablishmentID,
			TransactionType: enums.Payment,
			Amount:          amount,
			Description:     description,
			TransactionDate: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("error creating payment transaction: %w", err)
		}

		// 4. Update the credit account balance
		creditAccount.CurrentBalance -= amount
		if err := tx.Save(&creditAccount).Error; err != nil {
			return fmt.Errorf("error updating credit account balance: %w", err)
		}

		// 5. Create a CreditAccountHistory record
		historyEntry := entities.CreditAccountHistory{
			CreditAccountID: creditAccountID,
			TransactionDate: time.Now(),
			TransactionType: enums.Payment,
			Amount:          -amount, // Payment reduces balance
			Balance:         creditAccount.CurrentBalance,
			Description:     description,
		}
		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creating credit account history: %w", err)
		}

		// 6. Unblock the account if it was blocked and balance is 0 or less
		if creditAccount.IsBlocked && creditAccount.CurrentBalance <= 0 {
			creditAccount.IsBlocked = false
			if err := tx.Save(&creditAccount).Error; err != nil {
				return fmt.Errorf("error unblocking credit account: %w", err)
			}
		}

		return nil
	})
}

func (r *creditAccountRepository) ApproveCreditRequest(creditRequest *entities.CreditRequest) (*response.CreditAccountResponse, error) {
	var creditAccountResponse *response.CreditAccountResponse // Variable to store the response

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create the CreditAccount entity from the request
		creditAccount := entities.CreditAccount{ // Create CreditAccount entity
			EstablishmentID:         creditRequest.EstablishmentID,
			ClientID:                creditRequest.ClientID,
			CreditLimit:             creditRequest.RequestedCreditLimit,
			MonthlyDueDate:          creditRequest.MonthlyDueDate,
			InterestRate:            creditRequest.InterestRate,
			InterestType:            creditRequest.InterestType,
			CreditType:              creditRequest.CreditType,
			GracePeriod:             creditRequest.GracePeriod,
			IsBlocked:               false,
			LastInterestAccrualDate: time.Now(),
			CurrentBalance:          0, // Initial balance is 0
			LateFeeRuleID:           creditRequest.Establishment.LateFeeRuleID,
		}

		// 2. Save the credit account
		if err := tx.Create(&creditAccount).Error; err != nil {
			return fmt.Errorf("error creating credit account: %w", err)
		}

		// 3. Update the credit request status
		now := time.Now()
		creditRequest.Status = entities.Approved
		creditRequest.ApprovedAt = &now
		if err := tx.Save(creditRequest).Error; err != nil { // No need to call a separate method
			return fmt.Errorf("error updating credit request status: %w", err)
		}

		return nil // Successful transaction
	})

	if err != nil {
		return nil, err // Return the error from the transaction
	}

	return creditAccountResponse, nil // Return the response object
}

func (r *creditAccountRepository) GetCreditRequestByID(id uint) (*entities.CreditRequest, error) {
	var creditRequest entities.CreditRequest
	if err := r.db.Preload("Establishment.Admin.User").Preload("Establishment.LateFeeRule").First(&creditRequest, id).Error; err != nil {
		return nil, err
	}
	return &creditRequest, nil
}

func (r *creditAccountRepository) GetPendingCreditRequests(establishmentID uint) ([]entities.CreditRequest, error) {
	var creditRequests []entities.CreditRequest
	if err := r.db.Where("establishment_id = ? AND status = ?", establishmentID, enums.Pending).Find(&creditRequests).Error; err != nil {
		return nil, fmt.Errorf("error retrieving credit requests: %w", err)
	}
	return creditRequests, nil
}

func (r *creditAccountRepository) AssignCreditAccountToClient(creditAccountID, clientID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var creditAccount entities.CreditAccount
		if err := tx.First(&creditAccount, creditAccountID).Error; err != nil {
			return fmt.Errorf("error retrieving credit account: %w", err)
		}

		// Check if the credit account is already assigned to a client
		if creditAccount.ClientID != 0 {
			return fmt.Errorf("credit account is already assigned to client %d", creditAccount.ClientID)
		}

		// Assign the client to the credit account
		creditAccount.ClientID = clientID
		if err := tx.Save(&creditAccount).Error; err != nil {
			return fmt.Errorf("error assigning credit account to client: %w", err)
		}

		return nil
	})
}

func (r *creditAccountRepository) GetOverdueBalance(userID uint) (float64, error) {
	var creditAccount entities.CreditAccount
	err := r.db.Joins("Client").Where("client.user_id = ?", userID).First(&creditAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // No hay cuenta de crédito, por lo tanto, no hay saldo vencido
		}
		return 0, fmt.Errorf("error al buscar la cuenta de crédito: %w", err)
	}

	// Verificar si la cuenta está vencida
	if !isAccountOverdue(creditAccount) {
		return 0, nil // La cuenta no está vencida
	}

	// Calcular el saldo vencido (puedes tener una lógica más compleja aquí)
	overdueBalance := creditAccount.CurrentBalance

	return overdueBalance, nil
}

// isAccountOverdue verifica si la cuenta está vencida
func isAccountOverdue(creditAccount entities.CreditAccount) bool {
	today := time.Now()
	dueDate := time.Date(today.Year(), today.Month(), creditAccount.MonthlyDueDate, 0, 0, 0, 0, time.UTC)

	return today.After(dueDate) && creditAccount.CurrentBalance > 0
}

func (r *creditAccountRepository) GetCreditAccountByUserID(userID uint) (*entities.CreditAccount, error) {
	var creditAccount entities.CreditAccount

	// Une la tabla CreditAccount con la tabla Client para filtrar por UserID
	err := r.db.Joins("JOIN clients ON credit_accounts.client_id = clients.id").Where("clients.user_id = ?", userID).First(&creditAccount).Error
	if err != nil {
		return nil, err
	}

	return &creditAccount, nil
}
