package inventory

import (
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeIn       MovementType = "in"
	MovementTypeOut      MovementType = "out"
	MovementTypeTransfer MovementType = "transfer"
)

// ReferenceType represents the source of a stock movement
type ReferenceType string

const (
	ReferenceTypePurchaseOrder ReferenceType = "purchase_order"
	ReferenceTypeOrder        ReferenceType = "order"
	ReferenceTypeAdjustment   ReferenceType = "adjustment"
	ReferenceTypeReturn       ReferenceType = "return"
	ReferenceTypeTransfer     ReferenceType = "transfer"
)

// StockMovement represents a change in inventory stock
type StockMovement struct {
	common.Entity
	MovementID    uint         `json:"movement_id"`
	ProductID     uint         `json:"product_id"`
	VariantID     *uint        `json:"variant_id,omitempty"`
	WarehouseID   uint         `json:"warehouse_id"`
	Quantity      int          `json:"quantity"`
	Type          MovementType  `json:"type"`
	ReferenceType ReferenceType `json:"reference_type"`
	ReferenceID   uint         `json:"reference_id"`
	Notes         string       `json:"notes"`
	CreatedAt     time.Time    `json:"created_at"`
	StaffID       uint         `json:"staff_id"`
	Warehouse     *Warehouse   `json:"warehouse,omitempty"`
}
