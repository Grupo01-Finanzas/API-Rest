package entities

type CreditRequestStatus string

const (
	Pending  CreditRequestStatus = "PENDING"
	Approved CreditRequestStatus = "APPROVED"
	Rejected CreditRequestStatus = "REJECTED"
)
