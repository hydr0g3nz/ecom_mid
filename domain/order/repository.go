package order

import (
	"time"
)

// OrderRepository defines the interface for order operations
type OrderRepository interface {
	FindByID(id uint) (*Order, error)
	FindByOrderNumber(orderNumber string) (*Order, error)
	FindByCustomer(customerID uint, page, limit int) ([]*Order, error)
	FindByStatus(status OrderStatus, page, limit int) ([]*Order, error)
	FindByDateRange(startDate, endDate time.Time, page, limit int) ([]*Order, error)
	FindWithFilter(filter map[string]interface{}, page, limit int) ([]*Order, error)
	Create(order *Order) error
	Update(order *Order) error
	UpdateStatus(orderID uint, status OrderStatus, comment string, staffID *uint) error
	AddItem(item *OrderItem) error
	UpdateItem(item *OrderItem) error
	RemoveItem(itemID uint) error
	FindOrderItems(orderID uint) ([]*OrderItem, error)
	FindOrderHistory(orderID uint) ([]*OrderStatusHistory, error)
}

// PaymentMethodRepository defines the interface for payment method operations
type PaymentMethodRepository interface {
	FindByID(id uint) (*PaymentMethod, error)
	FindActive() ([]*PaymentMethod, error)
	FindAll() ([]*PaymentMethod, error)
	Create(paymentMethod *PaymentMethod) error
	Update(paymentMethod *PaymentMethod) error
	Delete(id uint) error
}

// TransactionRepository defines the interface for transaction operations
type TransactionRepository interface {
	FindByID(id uint) (*Transaction, error)
	FindByOrder(orderID uint) ([]*Transaction, error)
	FindByReferenceNumber(refNumber string) (*Transaction, error)
	FindByStatus(status PaymentStatus) ([]*Transaction, error)
	FindByDateRange(startDate, endDate time.Time, page, limit int) ([]*Transaction, error)
	Create(transaction *Transaction) error
	Update(transaction *Transaction) error
}

// RefundRepository defines the interface for refund operations
type RefundRepository interface {
	FindByID(id uint) (*Refund, error)
	FindByOrder(orderID uint) ([]*Refund, error)
	FindByTransaction(transactionID uint) ([]*Refund, error)
	FindByStatus(status string) ([]*Refund, error)
	Create(refund *Refund) error
	Update(refund *Refund) error
	ProcessRefund(refundID uint, status string, processedBy uint, notes string) error
}

// ShipmentRepository defines the interface for shipment operations
type ShipmentRepository interface {
	FindByID(id uint) (*Shipment, error)
	FindByOrder(orderID uint) ([]*Shipment, error)
	FindByTrackingNumber(trackingNumber string) (*Shipment, error)
	FindByStatus(status ShipmentStatus) ([]*Shipment, error)
	Create(shipment *Shipment) error
	Update(shipment *Shipment) error
	UpdateStatus(shipmentID uint, status ShipmentStatus) error
}

// DocumentRepository defines the interface for document operations
type DocumentRepository interface {
	FindByID(id uint) (*Document, error)
	FindByReference(referenceType string, referenceID uint) ([]*Document, error)
	FindByType(documentType DocumentType) ([]*Document, error)
	FindByDocumentNumber(documentNumber string) (*Document, error)
	Create(document *Document) error
	Update(document *Document) error
	GenerateDocument(documentType DocumentType, referenceType string, referenceID uint, createdBy uint) (*Document, error)
}

// DocumentTemplateRepository defines the interface for document template operations
type DocumentTemplateRepository interface {
	FindByID(id uint) (*DocumentTemplate, error)
	FindByType(documentType DocumentType) ([]*DocumentTemplate, error)
	FindDefault(documentType DocumentType) (*DocumentTemplate, error)
	Create(template *DocumentTemplate) error
	Update(template *DocumentTemplate) error
	Delete(id uint) error
	SetDefault(templateID uint) error
}
