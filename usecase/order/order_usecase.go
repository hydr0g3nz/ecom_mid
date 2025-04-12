package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
	"github.com/hydr0g3nz/ecom_mid/domain/order"
	"github.com/hydr0g3nz/ecom_mid/domain/product"
	"github.com/hydr0g3nz/ecom_mid/domain/user"
	"github.com/hydr0g3nz/ecom_mid/usecase/inventory"
)

// OrderUseCase contains the business logic for order operations
type OrderUseCase struct {
	orderRepo          order.OrderRepository
	customerRepo       user.CustomerRepository
	productRepo        product.ProductRepository
	variantRepo        product.ProductVariantRepository
	paymentMethodRepo  order.PaymentMethodRepository
	transactionRepo    order.TransactionRepository
	shipmentRepo       order.ShipmentRepository
	documentRepo       order.DocumentRepository
	inventoryUseCase   *inventory.InventoryUseCase
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(
	orderRepo order.OrderRepository,
	customerRepo user.CustomerRepository,
	productRepo product.ProductRepository,
	variantRepo product.ProductVariantRepository,
	paymentMethodRepo order.PaymentMethodRepository,
	transactionRepo order.TransactionRepository,
	shipmentRepo order.ShipmentRepository,
	documentRepo order.DocumentRepository,
	inventoryUseCase *inventory.InventoryUseCase,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:         orderRepo,
		customerRepo:      customerRepo,
		productRepo:       productRepo,
		variantRepo:       variantRepo,
		paymentMethodRepo: paymentMethodRepo,
		transactionRepo:   transactionRepo,
		shipmentRepo:      shipmentRepo,
		documentRepo:      documentRepo,
		inventoryUseCase:  inventoryUseCase,
	}
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(
	customerID uint,
	paymentMethodID uint,
	shippingAddressID uint,
	billingAddressID uint,
	notes string,
) (*order.Order, error) {
	// Verify customer exists
	customer, err := uc.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, err
	}
	
	if customer == nil {
		return nil, errors.New("customer not found")
	}
	
	// Verify payment method exists and is active
	paymentMethod, err := uc.paymentMethodRepo.FindByID(paymentMethodID)
	if err != nil {
		return nil, err
	}
	
	if paymentMethod == nil || !paymentMethod.IsActive {
		return nil, errors.New("payment method not found or inactive")
	}
	
	// Verify addresses exist and belong to customer
	addresses, err := uc.customerRepo.FindAddressesByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	
	shippingAddressFound := false
	billingAddressFound := false
	
	for _, addr := range addresses {
		if addr.AddressID == shippingAddressID {
			shippingAddressFound = true
		}
		if addr.AddressID == billingAddressID {
			billingAddressFound = true
		}
	}
	
	if !shippingAddressFound {
		return nil, errors.New("shipping address not found for this customer")
	}
	
	if !billingAddressFound {
		return nil, errors.New("billing address not found for this customer")
	}
	
	// Generate unique order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), customerID)
	
	// Create new order
	newOrder, err := order.NewOrder(
		customerID,
		orderNumber,
		paymentMethodID,
		shippingAddressID,
		billingAddressID,
		notes,
	)
	if err != nil {
		return nil, err
	}
	
	// Save order to repository
	err = uc.orderRepo.Create(newOrder)
	if err != nil {
		return nil, err
	}
	
	return newOrder, nil
}

// AddItemToOrder adds a product to an order
func (uc *OrderUseCase) AddItemToOrder(
	orderID uint,
	productID uint,
	variantID *uint,
	quantity int,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	// Find order
	ord, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	if ord == nil {
		return errors.New("order not found")
	}
	
	// Verify order is in pending status
	if ord.Status != order.OrderStatusPending {
		return errors.New("cannot add items to a non-pending order")
	}
	
	// Find product
	prod, err := uc.productRepo.FindByID(productID)
	if err != nil {
		return err
	}
	
	if prod == nil {
		return errors.New("product not found")
	}
	
	// Determine which product/variant to use
	sku := prod.SKU
	name := prod.Name
	price := prod.GetCurrentPrice()
	
	if variantID != nil {
		// If variant specified, find it
		variantFound := false
		for _, variant := range prod.Variants {
			if variant.VariantID == *variantID {
				sku = variant.SKU
				price = variant.GetCurrentPrice()
				variantFound = true
				break
			}
		}
		
		if !variantFound {
			return errors.New("variant not found for this product")
		}
	}
	
	// Calculate item totals
	subtotal, err := price.Multiply(float64(quantity))
	if err != nil {
		return err
	}
	
	// Initialize with zero for tax and discount
	tax, _ := vo.NewMoney(0, price.Currency)
	discount, _ := vo.NewMoney(0, price.Currency)
	
	// Create order item
	item := order.OrderItem{
		OrderID:   orderID,
		ProductID: productID,
		VariantID: variantID,
		SKU:       sku,
		Name:      name,
		Quantity:  quantity,
		UnitPrice: price,
		Subtotal:  subtotal,
		Tax:       tax,
		Discount:  discount,
		Total:     subtotal, // Without tax and discount for now
	}
	
	// Add item to order
	err = ord.AddItem(item)
	if err != nil {
		return err
	}
	
	// Reserve stock in inventory
	// Note: In a real implementation, you would likely handle this differently,
	// possibly with a default warehouse selection strategy
	// For simplicity, we're skipping the actual inventory reservation here
	
	// Save updated order
	return uc.orderRepo.Update(ord)
}

// ProcessPayment processes a payment for an order
func (uc *OrderUseCase) ProcessPayment(
	orderID uint,
	referenceNumber string,
	gatewayResponse string,
	gatewayTransactionID string,
) error {
	// Find order
	ord, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	if ord == nil {
		return errors.New("order not found")
	}
	
	// Create transaction
	transaction := order.Transaction{
		OrderID:              orderID,
		PaymentMethodID:      ord.PaymentMethodID,
		TransactionDate:      time.Now(),
		Amount:               ord.TotalAmount,
		Status:               order.PaymentStatusPaid,
		ReferenceNumber:      referenceNumber,
		GatewayResponse:      gatewayResponse,
		GatewayTransactionID: gatewayTransactionID,
	}
	
	// Save transaction
	err = uc.transactionRepo.Create(&transaction)
	if err != nil {
		return err
	}
	
	// Update order payment status
	ord.UpdatePaymentStatus(order.PaymentStatusPaid)
	
	// Update order status to processing
	err = ord.UpdateStatus(order.OrderStatusProcessing, "Payment received, order processing", nil)
	if err != nil {
		return err
	}
	
	// Save updated order
	return uc.orderRepo.Update(ord)
}

// CreateShipment creates a shipment for an order
func (uc *OrderUseCase) CreateShipment(
	orderID uint,
	trackingNumber string,
	carrier string,
	expectedDeliveryDate *time.Time,
) error {
	// Find order
	ord, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	if ord == nil {
		return errors.New("order not found")
	}
	
	// Verify order is in processing status
	if ord.Status != order.OrderStatusProcessing {
		return errors.New("can only create shipment for processing orders")
	}
	
	// Create shipment
	shipment := order.Shipment{
		OrderID:              orderID,
		TrackingNumber:       trackingNumber,
		Carrier:              carrier,
		ExpectedDeliveryDate: expectedDeliveryDate,
		Status:               order.ShipmentStatusPending,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	
	// Add shipment to order
	err = ord.AddShipment(shipment)
	if err != nil {
		return err
	}
	
	// Create shipment record
	err = uc.shipmentRepo.Create(&shipment)
	if err != nil {
		return err
	}
	
	// Commit reserved inventory (convert to actual deduction)
	// Note: In a real implementation, you would handle this with proper warehouse selection
	// For simplicity, we're skipping the actual inventory commitment here
	
	// Mark shipment as shipped
	shipment.MarkAsShipped()
	
	// Update shipment status
	err = uc.shipmentRepo.Update(&shipment)
	if err != nil {
		return err
	}
	
	// Save updated order
	return uc.orderRepo.Update(ord)
}

// MarkOrderDelivered marks an order as delivered
func (uc *OrderUseCase) MarkOrderDelivered(orderID uint, staffID *uint) error {
	// Find order
	ord, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	if ord == nil {
		return errors.New("order not found")
	}
	
	// Verify order is in shipped status
	if ord.Status != order.OrderStatusShipped {
		return errors.New("can only mark shipped orders as delivered")
	}
	
	// Update shipments status
	for _, shipment := range ord.Shipments {
		shipment.MarkAsDelivered()
		
		err = uc.shipmentRepo.Update(&shipment)
		if err != nil {
			return err
		}
	}
	
	// Update order status
	err = ord.UpdateStatus(order.OrderStatusDelivered, "Order delivered", staffID)
	if err != nil {
		return err
	}
	
	// Award loyalty points to customer (if applicable)
	// This would be implemented based on business rules
	
	// Save updated order
	return uc.orderRepo.Update(ord)
}

// CancelOrder cancels an order
func (uc *OrderUseCase) CancelOrder(orderID uint, reason string, staffID *uint) error {
	// Find order
	ord, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	if ord == nil {
		return errors.New("order not found")
	}
	
	// Cancel order
	err = ord.Cancel(reason, staffID)
	if err != nil {
		return err
	}
	
	// Release reserved inventory
	// Note: In a real implementation, you would handle this with proper warehouse selection
	// For simplicity, we're skipping the actual inventory release here
	
	// Save updated order
	return uc.orderRepo.Update(ord)
}

// GenerateInvoice generates an invoice for an order
func (uc *OrderUseCase) GenerateInvoice(orderID uint, staffID uint) (*order.Document, error) {
	// Find order
	ord, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	
	if ord == nil {
		return nil, errors.New("order not found")
	}
	
	// Generate invoice document
	document, err := uc.documentRepo.GenerateDocument(order.DocumentTypeInvoice, "order", orderID, staffID)
	if err != nil {
		return nil, err
	}
	
	return document, nil
}

// GetOrderByID gets an order by ID
func (uc *OrderUseCase) GetOrderByID(orderID uint) (*order.Order, error) {
	return uc.orderRepo.FindByID(orderID)
}

// GetOrderByNumber gets an order by order number
func (uc *OrderUseCase) GetOrderByNumber(orderNumber string) (*order.Order, error) {
	return uc.orderRepo.FindByOrderNumber(orderNumber)
}

// GetOrdersByCustomer gets orders for a customer
func (uc *OrderUseCase) GetOrdersByCustomer(customerID uint, page, limit int) ([]*order.Order, error) {
	return uc.orderRepo.FindByCustomer(customerID, page, limit)
}

// GetOrdersByStatus gets orders by status
func (uc *OrderUseCase) GetOrdersByStatus(status string, page, limit int) ([]*order.Order, error) {
	var orderStatus order.OrderStatus
	switch status {
	case "pending":
		orderStatus = order.OrderStatusPending
	case "processing":
		orderStatus = order.OrderStatusProcessing
	case "shipped":
		orderStatus = order.OrderStatusShipped
	case "delivered":
		orderStatus = order.OrderStatusDelivered
	case "cancelled":
		orderStatus = order.OrderStatusCancelled
	default:
		return nil, errors.New("invalid order status")
	}
	
	return uc.orderRepo.FindByStatus(orderStatus, page, limit)
}
