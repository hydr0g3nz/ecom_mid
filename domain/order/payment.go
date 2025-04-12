package order

import (
	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
	"time"
)

// PaymentMethodFeesType represents the type of fees for a payment method
type PaymentMethodFeesType string

const (
	FeesTypeFixed      PaymentMethodFeesType = "fixed"
	FeesTypePercentage PaymentMethodFeesType = "percentage"
)

// PaymentMethod represents a payment method available in the system
type PaymentMethod struct {
	common.Entity
	PaymentMethodID uint                `json:"payment_method_id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	IsActive        bool                `json:"is_active"`
	FeesType        PaymentMethodFeesType `json:"fees_type"`
	FeesAmount      float64             `json:"fees_amount"`
	SortOrder       int                 `json:"sort_order"`
}

// Transaction represents a payment transaction for an order
type Transaction struct {
	common.Entity
	TransactionID        uint          `json:"transaction_id"`
	OrderID              uint          `json:"order_id"`
	PaymentMethodID      uint          `json:"payment_method_id"`
	TransactionDate      time.Time     `json:"transaction_date"`
	Amount               vo.Money      `json:"amount"`
	Status               PaymentStatus `json:"status"`
	ReferenceNumber      string        `json:"reference_number"`
	GatewayResponse      string        `json:"gateway_response"`
	GatewayTransactionID string        `json:"gateway_transaction_id"`
	Order                *Order        `json:"order,omitempty"`
	PaymentMethod        *PaymentMethod `json:"payment_method,omitempty"`
}

// Refund represents a refund for an order
type Refund struct {
	common.Entity
	RefundID      uint          `json:"refund_id"`
	OrderID       uint          `json:"order_id"`
	TransactionID uint          `json:"transaction_id"`
	Amount        vo.Money      `json:"amount"`
	Reason        string        `json:"reason"`
	Status        string        `json:"status"`
	RefundDate    time.Time     `json:"refund_date"`
	ProcessedBy   *uint         `json:"processed_by,omitempty"`
	Notes         string        `json:"notes"`
	Order         *Order        `json:"order,omitempty"`
	Transaction   *Transaction  `json:"transaction,omitempty"`
}

// CalculatePaymentFee calculates additional fees for using a payment method
func (pm *PaymentMethod) CalculatePaymentFee(amount float64) float64 {
	if !pm.IsActive || pm.FeesAmount <= 0 {
		return 0
	}
	
	if pm.FeesType == FeesTypeFixed {
		return pm.FeesAmount
	}
	
	// If percentage-based
	return (amount * pm.FeesAmount) / 100
}
