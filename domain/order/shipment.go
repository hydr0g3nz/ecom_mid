package order

import (
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "pending"
	ShipmentStatusShipped   ShipmentStatus = "shipped"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusFailed    ShipmentStatus = "failed"
)

// Shipment represents a physical shipment of an order
type Shipment struct {
	common.Entity
	ShipmentID           uint           `json:"shipment_id"`
	OrderID              uint           `json:"order_id"`
	TrackingNumber       string         `json:"tracking_number"`
	Carrier              string         `json:"carrier"`
	ShippingDate         *time.Time     `json:"shipping_date,omitempty"`
	ExpectedDeliveryDate *time.Time     `json:"expected_delivery_date,omitempty"`
	Status               ShipmentStatus `json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	Order                *Order        `json:"order,omitempty"`
}

// MarkAsShipped updates the shipment status to shipped
func (s *Shipment) MarkAsShipped() {
	now := time.Now()
	s.Status = ShipmentStatusShipped
	s.ShippingDate = &now
	s.UpdatedAt = now
}

// MarkAsDelivered updates the shipment status to delivered
func (s *Shipment) MarkAsDelivered() {
	s.Status = ShipmentStatusDelivered
	s.UpdatedAt = time.Now()
}

// MarkAsFailed updates the shipment status to failed
func (s *Shipment) MarkAsFailed() {
	s.Status = ShipmentStatusFailed
	s.UpdatedAt = time.Now()
}
