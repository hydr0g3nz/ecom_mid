package inventory

import (
	"errors"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// Inventory represents the stock level of a product in a warehouse
type Inventory struct {
	common.Entity
	InventoryID      uint      `json:"inventory_id"`
	ProductID        uint      `json:"product_id"`
	VariantID        *uint     `json:"variant_id,omitempty"`
	WarehouseID      uint      `json:"warehouse_id"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	UpdatedAt        time.Time `json:"updated_at"`
	Warehouse        *Warehouse `json:"warehouse,omitempty"`
}

// GetAvailableQuantity returns the available quantity (total - reserved)
func (i *Inventory) GetAvailableQuantity() int {
	return i.Quantity - i.ReservedQuantity
}

// IsInStock checks if there is available stock
func (i *Inventory) IsInStock() bool {
	return i.GetAvailableQuantity() > 0
}

// CanFulfill checks if there is enough available stock to fulfill a quantity
func (i *Inventory) CanFulfill(quantity int) bool {
	return i.GetAvailableQuantity() >= quantity
}

// AddStock adds stock to the inventory
func (i *Inventory) AddStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	i.Quantity += quantity
	i.UpdatedAt = time.Now()
	return nil
}

// RemoveStock removes stock from the inventory
func (i *Inventory) RemoveStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	if i.Quantity < quantity {
		return errors.New("insufficient stock")
	}
	
	i.Quantity -= quantity
	i.UpdatedAt = time.Now()
	return nil
}

// ReserveStock reserves stock for an order
func (i *Inventory) ReserveStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	if i.GetAvailableQuantity() < quantity {
		return errors.New("insufficient available stock")
	}
	
	i.ReservedQuantity += quantity
	i.UpdatedAt = time.Now()
	return nil
}

// ReleaseReservedStock releases previously reserved stock
func (i *Inventory) ReleaseReservedStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	if i.ReservedQuantity < quantity {
		return errors.New("trying to release more than reserved")
	}
	
	i.ReservedQuantity -= quantity
	i.UpdatedAt = time.Now()
	return nil
}

// CommitReservedStock converts reserved stock to actual deduction
func (i *Inventory) CommitReservedStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	if i.ReservedQuantity < quantity {
		return errors.New("trying to commit more than reserved")
	}
	
	i.Quantity -= quantity
	i.ReservedQuantity -= quantity
	i.UpdatedAt = time.Now()
	return nil
}
