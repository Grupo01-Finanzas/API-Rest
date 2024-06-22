package enums

type PaymentStatus string

const (
	PENDING PaymentStatus = "PENDING"
	SUCCESS PaymentStatus = "SUCCESS"
	FAILED  PaymentStatus = "FAILED"
)