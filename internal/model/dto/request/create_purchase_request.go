package request

import "ApiRestFinance/internal/model/entities/enums"

// CreatePurchaseRequest holds the data to create a purchase
type CreatePurchaseRequest struct {
	EstablishmentID uint             `json:"establishment_id" binding:"required"`
	ProductIDs      []uint           `json:"product_ids" binding:"required"`
	CreditType      enums.CreditType `json:"credit_type" binding:"required"`
	Amount          float64          `json:"amount" binding:"required"`
}
