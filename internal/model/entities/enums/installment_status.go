package enums

type InstallmentStatus string

const (
	Pending InstallmentStatus = "PENDING"
	Paid    InstallmentStatus = "PAID"
	Overdue InstallmentStatus = "OVERDUE"
)
