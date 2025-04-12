package order

import (
	"errors"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
	"github.com/hydr0g3nz/ecom_mid/domain/user"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Order represents a customer's order
type Order struct {
	common.Entity
	OrderID          uint          `json:"order_id"`
	CustomerID       uint          `json:"customer_id"`
	OrderNumber      string        `json:"order_number"`
	OrderDate        time.Time     `json:"order_date"`
	Status           OrderStatus   `json:"status"`
	Subtotal         vo.Money      `json:"subtotal"`
	ShippingFee      vo.Money      `json:"shipping_fee"`
	TaxAmount        vo.Money      `json:"tax_amount"`
	DiscountAmount   vo.Money      `json:"discount_amount"`
	TotalAmount      vo.Money      `json:"total_amount"`
	PaymentMethodID  uint          `json:"payment_method_id"`
	PaymentStatus    PaymentStatus `json:"payment_status"`
	ShippingAddressID uint         `json:"shipping_address_id"`
	BillingAddressID  uint         `json:"billing_address_id"`
	Notes             string       `json:"notes"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	Customer          *user.Customer `json:"customer,omitempty"`
	PaymentMethod     *PaymentMethod `json:"payment_method,omitempty"`
	ShippingAddress   *user.CustomerAddress `json:"shipping_address,omitempty"`
	BillingAddress    *user.CustomerAddress `json:"billing_address,omitempty"`
	Items             []OrderItem `json:"items,omitempty"`
	StatusHistory     []OrderStatusHistory `json:"status_history,omitempty"`
	Shipments         []Shipment `json:"shipments,omitempty"`
	Transactions      []Transaction `json:"transactions,omitempty"`
}

// OrderItem represents a product in an order
type OrderItem struct {
	common.Entity
	OrderItemID uint    `json:"order_item_id"`
	OrderID     uint    `json:"order_id"`
	ProductID   uint    `json:"product_id"`
	VariantID   *uint   `json:"variant_id,omitempty"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   vo.Money `json:"unit_price"`
	Subtotal    vo.Money `json:"subtotal"`
	Tax         vo.Money `json:"tax"`
	Discount    vo.Money `json:"discount"`
	Total       vo.Money `json:"total"`
}

// OrderStatusHistory represents a change in order status
type OrderStatusHistory struct {
	common.Entity
	HistoryID uint        `json:"history_id"`
	OrderID   uint        `json:"order_id"`
	Status    OrderStatus `json:"status"`
	Comment   string      `json:"comment"`
	StaffID   *uint       `json:"staff_id,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

// NewOrder creates a new order with validation
func NewOrder(
	customerID uint,
	orderNumber string,
	paymentMethodID uint,
	shippingAddressID uint,
	billingAddressID uint,
	notes string,
) (*Order, error) {
	if customerID == 0 {
		return nil, errors.New("customer ID is required")
	}
	
	if orderNumber == "" {
		return nil, errors.New("order number is required")
	}
	
	if paymentMethodID == 0 {
		return nil, errors.New("payment method is required")
	}
	
	if shippingAddressID == 0 {
		return nil, errors.New("shipping address is required")
	}
	
	if billingAddressID == 0 {
		return nil, errors.New("billing address is required")
	}
	
	// Initialize with zero for monetary values
	subTotal, _ := vo.NewMoney(0, "THB")
	shippingFee, _ := vo.NewMoney(0, "THB")
	taxAmount, _ := vo.NewMoney(0, "THB")
	discountAmount, _ := vo.NewMoney(0, "THB")
	totalAmount, _ := vo.NewMoney(0, "THB")
	
	order := &Order{
		CustomerID:       customerID,
		OrderNumber:      orderNumber,
		OrderDate:        time.Now(),
		Status:           OrderStatusPending,
		Subtotal:         subTotal,
		ShippingFee:      shippingFee,
		TaxAmount:        taxAmount,
		DiscountAmount:   discountAmount,
		TotalAmount:      totalAmount,
		PaymentMethodID:  paymentMethodID,
		PaymentStatus:    PaymentStatusPending,
		ShippingAddressID: shippingAddressID,
		BillingAddressID:  billingAddressID,
		Notes:             notes,
		Items:             []OrderItem{},
		StatusHistory:     []OrderStatusHistory{},
		Shipments:         []Shipment{},
		Transactions:      []Transaction{},
	}
	
	// Add initial status history
	order.AddStatusHistory(OrderStatusPending, "Order created", nil)
	
	return order, nil
}

// AddItem adds a product to the order
func (o *Order) AddItem(item OrderItem) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot add items to a non-pending order")
	}
	
	// Calculate item total
	total, err := item.UnitPrice.Multiply(float64(item.Quantity))
	if err != nil {
		return err
	}
	
	// Apply tax and discount calculations here if needed
	
	// Set the calculated values
	item.Subtotal = total
	item.Total = total
	
	// Add item to order
	o.Items = append(o.Items, item)
	
	// Recalculate order totals
	return o.recalculateOrderTotals()
}

// UpdateItem updates an existing item in the order
func (o *Order) UpdateItem(itemID uint, quantity int) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot update items in a non-pending order")
	}
	
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	// Find and update the item
	itemFound := false
	for i, item := range o.Items {
		if item.OrderItemID == itemID {
			// Calculate new subtotal
			newTotal, err := item.UnitPrice.Multiply(float64(quantity))
			if err != nil {
				return err
			}
			
			o.Items[i].Quantity = quantity
			o.Items[i].Subtotal = newTotal
			o.Items[i].Total = newTotal
			
			itemFound = true
			break
		}
	}
	
	if !itemFound {
		return errors.New("item not found in order")
	}
	
	// Recalculate order totals
	return o.recalculateOrderTotals()
}

// RemoveItem removes an item from the order
func (o *Order) RemoveItem(itemID uint) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot remove items from a non-pending order")
	}
	
	// Find and remove the item
	itemIndex := -1
	for i, item := range o.Items {
		if item.OrderItemID == itemID {
			itemIndex = i
			break
		}
	}
	
	if itemIndex == -1 {
		return errors.New("item not found in order")
	}
	
	// Remove item
	o.Items = append(o.Items[:itemIndex], o.Items[itemIndex+1:]...)
	
	// Recalculate order totals
	return o.recalculateOrderTotals()
}

// recalculateOrderTotals recalculates all order totals based on items
func (o *Order) recalculateOrderTotals() error {
	// Reset subtotal
	subtotal, _ := vo.NewMoney(0, "THB")
	
	// Sum all item totals
	for _, item := range o.Items {
		newSubtotal, err := subtotal.Add(item.Total)
		if err != nil {
			return err
		}
		subtotal = newSubtotal
	}
	
	o.Subtotal = subtotal
	
	// Calculate final total (subtotal + shipping - discounts + tax)
	total := o.Subtotal
	
	// Add shipping fee
	if !o.ShippingFee.IsZero() {
		newTotal, err := total.Add(o.ShippingFee)
		if err != nil {
			return err
		}
		total = newTotal
	}
	
	// Add tax
	if !o.TaxAmount.IsZero() {
		newTotal, err := total.Add(o.TaxAmount)
		if err != nil {
			return err
		}
		total = newTotal
	}
	
	// Subtract discount
	if !o.DiscountAmount.IsZero() {
		newTotal, err := total.Subtract(o.DiscountAmount)
		if err != nil {
			return err
		}
		total = newTotal
	}
	
	o.TotalAmount = total
	return nil
}

// SetShippingFee sets the shipping fee for the order
func (o *Order) SetShippingFee(fee vo.Money) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot update shipping fee in a non-pending order")
	}
	
	o.ShippingFee = fee
	return o.recalculateOrderTotals()
}

// SetDiscount sets a discount amount for the order
func (o *Order) SetDiscount(discount vo.Money) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot update discount in a non-pending order")
	}
	
	o.DiscountAmount = discount
	return o.recalculateOrderTotals()
}

// SetTaxAmount sets the tax amount for the order
func (o *Order) SetTaxAmount(tax vo.Money) error {
	if o.Status != OrderStatusPending {
		return errors.New("cannot update tax in a non-pending order")
	}
	
	o.TaxAmount = tax
	return o.recalculateOrderTotals()
}

// UpdateStatus updates the order status and adds to history
func (o *Order) UpdateStatus(status OrderStatus, comment string, staffID *uint) error {
	// Validate status transition
	if !isValidStatusTransition(o.Status, status) {
		return errors.New("invalid status transition")
	}
	
	// Update status
	o.Status = status
	
	// Add to history
	return o.AddStatusHistory(status, comment, staffID)
}

// AddStatusHistory adds a status change to history
func (o *Order) AddStatusHistory(status OrderStatus, comment string, staffID *uint) error {
	history := OrderStatusHistory{
		OrderID:   o.OrderID,
		Status:    status,
		Comment:   comment,
		StaffID:   staffID,
		CreatedAt: time.Now(),
	}
	
	o.StatusHistory = append(o.StatusHistory, history)
	return nil
}

// UpdatePaymentStatus updates the payment status of the order
func (o *Order) UpdatePaymentStatus(status PaymentStatus) {
	o.PaymentStatus = status
	
	// If order is paid and still pending, move to processing
	if status == PaymentStatusPaid && o.Status == OrderStatusPending {
		o.UpdateStatus(OrderStatusProcessing, "Payment received, order processing", nil)
	}
}

// AddShipment adds a shipment to the order
func (o *Order) AddShipment(shipment Shipment) error {
	if o.Status != OrderStatusProcessing {
		return errors.New("can only add shipment to processing orders")
	}
	
	shipment.OrderID = o.OrderID
	o.Shipments = append(o.Shipments, shipment)
	
	// Update order status to shipped
	return o.UpdateStatus(OrderStatusShipped, "Order shipped", nil)
}

// Cancel cancels the order
func (o *Order) Cancel(reason string, staffID *uint) error {
	if o.Status == OrderStatusDelivered {
		return errors.New("cannot cancel a delivered order")
	}
	
	return o.UpdateStatus(OrderStatusCancelled, reason, staffID)
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(current, new OrderStatus) bool {
	switch current {
	case OrderStatusPending:
		return new == OrderStatusProcessing || new == OrderStatusCancelled
	case OrderStatusProcessing:
		return new == OrderStatusShipped || new == OrderStatusCancelled
	case OrderStatusShipped:
		return new == OrderStatusDelivered || new == OrderStatusCancelled
	case OrderStatusDelivered:
		return false // Terminal state
	case OrderStatusCancelled:
		return false // Terminal state
	default:
		return false
	}
}
