package enums

type RecipientType string

const (
	RolClient        RecipientType = "CLIENT"
	RolEstablishment RecipientType = "ESTABLISHMENT"
	RolAdmin         RecipientType = "ADMIN"
)
