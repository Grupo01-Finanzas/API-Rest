package enums

type TransactionType string

const (
	Purchase            TransactionType = "PURCHASE"
	Payment             TransactionType = "PAYMENT"
	InterestAccrual     TransactionType = "INTEREST_ACCRUAL"
	LateFeeApplied      TransactionType = "LATE_FEE_APPLIED"
	CreditLimitIncrease TransactionType = "CREDIT_LIMIT_INCREASE"
	CreditLimitDecrease TransactionType = "CREDIT_LIMIT_DECREASE"
	EarlyPayment        TransactionType = "EARLY_PAYMENT"
	AccountBlocked      TransactionType = "ACCOUNT_BLOCKED"
	AccountUnblocked    TransactionType = "ACCOUNT_UNBLOCKED"
)
