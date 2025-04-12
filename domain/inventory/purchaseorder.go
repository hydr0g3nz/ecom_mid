package inventory

import (
	"errors"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
)

// PurchaseOrderStatus represents the status of a purchase order
type PurchaseOrderStatus string

const (
	POStatusDraft     PurchaseOrderStatus = "draft"
	POStatusPending   PurchaseOrderStatus = "pending"
	POStatusConfirmed PurchaseOrderStatus = "confirmed"
	POStatusReceived  PurchaseOrderStatus = "received"
	POStatusCancelled PurchaseOrderStatus = "cancelled"
)

// PurchaseOrder represents an order to a supplier for products
type PurchaseOrder struct {
	common.Entity
	POID           uint                `json:"po_id"`
	SupplierID     uint                `json:"supplier_id"`
	OrderDate      time.Time           `json:"order_date"`
	ExpectedDate   *time.Time          `json:"expected_date,omitempty"`
	Status         PurchaseOrderStatus `json:"status"`
	Notes          string              `json:"notes"`
	TotalAmount    vo.Money            `json:"total_amount"`
	StaffID        uint                `json:"staff_id"`
	Supplier       *Supplier           `json:"supplier,omitempty"`
	Items          []PurchaseOrderItem `json:"items,omitempty"`
}

// PurchaseOrderItem represents a line item in a purchase order
type PurchaseOrderItem struct {
	common.Entity
	POItemID   uint     `json:"po_item_id"`
	POID       uint     `json:"po_id"`
	ProductID  uint     `json:"product_id"`
	VariantID  *uint    `json:"variant_id,omitempty"`
	Quantity   int      `json:"quantity"`
	UnitCost   vo.Money `json:"unit_cost"`
	TotalCost  vo.Money `json:"total_cost"`
}

// AddItem adds a product to the purchase order
func (po *PurchaseOrder) AddItem(item PurchaseOrderItem) error {
	if po.Status != POStatusDraft && po.Status != POStatusPending {
		return errors.New("cannot add items to a confirmed or received purchase order")
	}
	
	// Check if product already exists in PO
	for i, existingItem := range po.Items {
		// If same product and variant, just update quantity
		if existingItem.ProductID == item.ProductID && 
		   ((existingItem.VariantID == nil && item.VariantID == nil) || 
		    (existingItem.VariantID != nil && item.VariantID != nil && *existingItem.VariantID == *item.VariantID)) {
			
			// Update quantity
			newQuantity := existingItem.Quantity + item.Quantity
			
			// Recalculate total cost
			newTotalCost, err := existingItem.UnitCost.Multiply(float64(newQuantity))
			if err != nil {
				return err
			}
			
			po.Items[i].Quantity = newQuantity
			po.Items[i].TotalCost = newTotalCost
			
			// Recalculate PO total
			return po.recalculateTotal()
		}
	}
	
	// Add new item
	po.Items = append(po.Items, item)
	
	// Recalculate PO total
	return po.recalculateTotal()
}

// UpdateItem updates an existing item in the purchase order
func (po *PurchaseOrder) UpdateItem(itemID uint, quantity int, unitCost vo.Money) error {
	if po.Status != POStatusDraft && po.Status != POStatusPending {
		return errors.New("cannot update items in a confirmed or received purchase order")
	}
	
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	// Find and update the item
	itemFound := false
	for i, item := range po.Items {
		if item.POItemID == itemID {
			// Calculate new total cost
			newTotalCost, err := unitCost.Multiply(float64(quantity))
			if err != nil {
				return err
			}
			
			po.Items[i].Quantity = quantity
			po.Items[i].UnitCost = unitCost
			po.Items[i].TotalCost = newTotalCost
			
			itemFound = true
			break
		}
	}
	
	if !itemFound {
		return errors.New("item not found in purchase order")
	}
	
	// Recalculate PO total
	return po.recalculateTotal()
}

// RemoveItem removes an item from the purchase order
func (po *PurchaseOrder) RemoveItem(itemID uint) error {
	if po.Status != POStatusDraft && po.Status != POStatusPending {
		return errors.New("cannot remove items from a confirmed or received purchase order")
	}
	
	// Find and remove the item
	itemIndex := -1
	for i, item := range po.Items {
		if item.POItemID == itemID {
			itemIndex = i
			break
		}
	}
	
	if itemIndex == -1 {
		return errors.New("item not found in purchase order")
	}
	
	// Remove item
	po.Items = append(po.Items[:itemIndex], po.Items[itemIndex+1:]...)
	
	// Recalculate PO total
	return po.recalculateTotal()
}

// recalculateTotal recalculates the total amount of the purchase order
func (po *PurchaseOrder) recalculateTotal() error {
	if len(po.Items) == 0 {
		po.TotalAmount, _ = vo.NewMoney(0, "THB")
		return nil
	}
	
	// Initialize with first item's currency
	total := po.Items[0].TotalCost
	
	// Add all other items
	for i := 1; i < len(po.Items); i++ {
		newTotal, err := total.Add(po.Items[i].TotalCost)
		if err != nil {
			return err
		}
		total = newTotal
	}
	
	po.TotalAmount = total
	return nil
}

// ConfirmOrder changes the status to confirmed
func (po *PurchaseOrder) ConfirmOrder() error {
	if po.Status != POStatusPending {
		return errors.New("only pending purchase orders can be confirmed")
	}
	
	if len(po.Items) == 0 {
		return errors.New("cannot confirm empty purchase order")
	}
	
	po.Status = POStatusConfirmed
	return nil
}

// ReceiveOrder changes the status to received
func (po *PurchaseOrder) ReceiveOrder() error {
	if po.Status != POStatusConfirmed {
		return errors.New("only confirmed purchase orders can be received")
	}
	
	po.Status = POStatusReceived
	return nil
}

// CancelOrder changes the status to cancelled
func (po *PurchaseOrder) CancelOrder() error {
	if po.Status == POStatusReceived {
		return errors.New("received purchase orders cannot be cancelled")
	}
	
	po.Status = POStatusCancelled
	return nil
}
