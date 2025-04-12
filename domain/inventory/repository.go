package inventory

// WarehouseRepository defines the interface for warehouse operations
type WarehouseRepository interface {
	FindByID(id uint) (*Warehouse, error)
	FindByName(name string) (*Warehouse, error)
	FindAll() ([]*Warehouse, error)
	Create(warehouse *Warehouse) error
	Update(warehouse *Warehouse) error
	Delete(id uint) error
}

// InventoryRepository defines the interface for inventory operations
type InventoryRepository interface {
	FindByID(id uint) (*Inventory, error)
	FindByProductAndWarehouse(productID, warehouseID uint, variantID *uint) (*Inventory, error)
	FindByProduct(productID uint, variantID *uint) ([]*Inventory, error)
	FindByWarehouse(warehouseID uint) ([]*Inventory, error)
	FindLowStock(threshold int) ([]*Inventory, error)
	Create(inventory *Inventory) error
	Update(inventory *Inventory) error
	Delete(id uint) error
	AddStock(productID, warehouseID uint, variantID *uint, quantity int, referenceType string, referenceID uint, staffID uint, notes string) error
	RemoveStock(productID, warehouseID uint, variantID *uint, quantity int, referenceType string, referenceID uint, staffID uint, notes string) error
	ReserveStock(productID, warehouseID uint, variantID *uint, quantity int, orderID uint, staffID uint) error
	ReleaseReservedStock(productID, warehouseID uint, variantID *uint, quantity int, orderID uint, staffID uint) error
	CommitReservedStock(productID, warehouseID uint, variantID *uint, quantity int, orderID uint, staffID uint) error
}

// StockMovementRepository defines the interface for stock movement operations
type StockMovementRepository interface {
	FindByID(id uint) (*StockMovement, error)
	FindByProduct(productID uint, variantID *uint, page, limit int) ([]*StockMovement, error)
	FindByWarehouse(warehouseID uint, page, limit int) ([]*StockMovement, error)
	FindByReference(referenceType string, referenceID uint) ([]*StockMovement, error)
	FindByDateRange(startDate, endDate time.Time, page, limit int) ([]*StockMovement, error)
	Create(movement *StockMovement) error
}

// SupplierRepository defines the interface for supplier operations
type SupplierRepository interface {
	FindByID(id uint) (*Supplier, error)
	FindByName(name string) (*Supplier, error)
	FindAll() ([]*Supplier, error)
	Create(supplier *Supplier) error
	Update(supplier *Supplier) error
	Delete(id uint) error
}

// PurchaseOrderRepository defines the interface for purchase order operations
type PurchaseOrderRepository interface {
	FindByID(id uint) (*PurchaseOrder, error)
	FindBySupplier(supplierID uint, page, limit int) ([]*PurchaseOrder, error)
	FindByStatus(status PurchaseOrderStatus, page, limit int) ([]*PurchaseOrder, error)
	FindByDateRange(startDate, endDate time.Time, page, limit int) ([]*PurchaseOrder, error)
	Create(po *PurchaseOrder) error
	Update(po *PurchaseOrder) error
	Delete(id uint) error
	AddItem(item *PurchaseOrderItem) error
	UpdateItem(item *PurchaseOrderItem) error
	RemoveItem(itemID uint) error
	ConfirmOrder(poID uint, staffID uint) error
	ReceiveOrder(poID uint, staffID uint) error
	CancelOrder(poID uint, staffID uint) error
}
