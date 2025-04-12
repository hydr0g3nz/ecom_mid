-- 1. ระบบสตาฟผู้ดูแล
CREATE TABLE Staff (
    staff_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    position VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active'
);

CREATE TABLE Role (
    role_id INT AUTO_INCREMENT PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE StaffRole (
    staff_role_id INT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT NOT NULL,
    role_id INT NOT NULL,
    FOREIGN KEY (staff_id) REFERENCES Staff(staff_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Role(role_id) ON DELETE CASCADE,
    UNIQUE KEY (staff_id, role_id)
);

CREATE TABLE Permission (
    permission_id INT AUTO_INCREMENT PRIMARY KEY,
    permission_name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE RolePermission (
    role_permission_id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES Role(role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES Permission(permission_id) ON DELETE CASCADE,
    UNIQUE KEY (role_id, permission_id)
);

-- 2. ระบบสมาชิก
CREATE TABLE Customer (
    customer_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    birth_date DATE,
    gender ENUM('male', 'female', 'other'),
    profile_image VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    points INT DEFAULT 0
);

CREATE TABLE CustomerAddress (
    address_id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    address_type ENUM('shipping', 'billing') NOT NULL,
    receiver_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    district VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (customer_id) REFERENCES Customer(customer_id) ON DELETE CASCADE
);

-- 3. ระบบ PIM (Product Information Management)
CREATE TABLE Category (
    category_id INT AUTO_INCREMENT PRIMARY KEY,
    parent_category_id INT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image VARCHAR(255),
    status ENUM('active', 'inactive') DEFAULT 'active',
    display_order INT DEFAULT 0,
    FOREIGN KEY (parent_category_id) REFERENCES Category(category_id) ON DELETE SET NULL
);

CREATE TABLE Product (
    product_id INT AUTO_INCREMENT PRIMARY KEY,
    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    short_description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    special_price DECIMAL(10, 2),
    special_price_start_date DATETIME,
    special_price_end_date DATETIME,
    cost DECIMAL(10, 2),
    weight DECIMAL(10, 2),
    length DECIMAL(10, 2),
    width DECIMAL(10, 2),
    height DECIMAL(10, 2),
    status ENUM('active', 'inactive', 'draft') DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords VARCHAR(255)
);

CREATE TABLE ProductCategory (
    product_category_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    category_id INT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES Category(category_id) ON DELETE CASCADE,
    UNIQUE KEY (product_id, category_id)
);

CREATE TABLE AttributeGroup (
    attribute_group_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE Attribute (
    attribute_id INT AUTO_INCREMENT PRIMARY KEY,
    attribute_group_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    type ENUM('text', 'select', 'boolean', 'number', 'date') NOT NULL,
    is_filterable BOOLEAN DEFAULT FALSE,
    is_required BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (attribute_group_id) REFERENCES AttributeGroup(attribute_group_id) ON DELETE CASCADE,
    UNIQUE KEY (attribute_group_id, name)
);

CREATE TABLE AttributeOption (
    option_id INT AUTO_INCREMENT PRIMARY KEY,
    attribute_id INT NOT NULL,
    value VARCHAR(255) NOT NULL,
    FOREIGN KEY (attribute_id) REFERENCES Attribute(attribute_id) ON DELETE CASCADE,
    UNIQUE KEY (attribute_id, value)
);

CREATE TABLE ProductAttribute (
    product_attribute_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    attribute_id INT NOT NULL,
    value TEXT,
    option_id INT,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES Attribute(attribute_id) ON DELETE CASCADE,
    FOREIGN KEY (option_id) REFERENCES AttributeOption(option_id) ON DELETE SET NULL,
    UNIQUE KEY (product_id, attribute_id)
);

CREATE TABLE ProductVariant (
    variant_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    sku VARCHAR(50) NOT NULL UNIQUE,
    price DECIMAL(10, 2) NOT NULL,
    special_price DECIMAL(10, 2),
    stock_quantity INT DEFAULT 0,
    weight DECIMAL(10, 2),
    status ENUM('active', 'inactive') DEFAULT 'active',
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE
);

CREATE TABLE VariantAttribute (
    variant_attribute_id INT AUTO_INCREMENT PRIMARY KEY,
    variant_id INT NOT NULL,
    attribute_id INT NOT NULL,
    option_id INT NOT NULL,
    FOREIGN KEY (variant_id) REFERENCES ProductVariant(variant_id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES Attribute(attribute_id) ON DELETE CASCADE,
    FOREIGN KEY (option_id) REFERENCES AttributeOption(option_id) ON DELETE CASCADE,
    UNIQUE KEY (variant_id, attribute_id)
);

CREATE TABLE ProductImage (
    image_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    variant_id INT,
    image_path VARCHAR(255) NOT NULL,
    alt_text VARCHAR(255),
    sort_order INT DEFAULT 0,
    is_main BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES ProductVariant(variant_id) ON DELETE CASCADE
);

-- 4. ระบบ Inventory
CREATE TABLE Warehouse (
    warehouse_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    status ENUM('active', 'inactive') DEFAULT 'active'
);

CREATE TABLE Inventory (
    inventory_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    variant_id INT,
    warehouse_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES ProductVariant(variant_id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES Warehouse(warehouse_id) ON DELETE CASCADE,
    UNIQUE KEY (product_id, variant_id, warehouse_id)
);

CREATE TABLE Supplier (
    supplier_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    contact_person VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    status ENUM('active', 'inactive') DEFAULT 'active'
);

CREATE TABLE StockMovement (
    movement_id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    variant_id INT,
    warehouse_id INT NOT NULL,
    quantity INT NOT NULL,
    type ENUM('in', 'out', 'transfer') NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    reference_id INT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    staff_id INT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES ProductVariant(variant_id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES Warehouse(warehouse_id) ON DELETE CASCADE,
    FOREIGN KEY (staff_id) REFERENCES Staff(staff_id) ON DELETE CASCADE
);

CREATE TABLE PurchaseOrder (
    po_id INT AUTO_INCREMENT PRIMARY KEY,
    supplier_id INT NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expected_date DATE,
    status ENUM('draft', 'pending', 'confirmed', 'received', 'cancelled') DEFAULT 'draft',
    notes TEXT,
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    staff_id INT NOT NULL,
    FOREIGN KEY (supplier_id) REFERENCES Supplier(supplier_id) ON DELETE CASCADE,
    FOREIGN KEY (staff_id) REFERENCES Staff(staff_id) ON DELETE CASCADE
);

CREATE TABLE PurchaseOrderItem (
    po_item_id INT AUTO_INCREMENT PRIMARY KEY,
    po_id INT NOT NULL,
    product_id INT NOT NULL,
    variant_id INT,
    quantity INT NOT NULL,
    unit_cost DECIMAL(10, 2) NOT NULL,
    total_cost DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (po_id) REFERENCES PurchaseOrder(po_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES ProductVariant(variant_id) ON DELETE CASCADE
);

-- 5. ระบบ Order
CREATE TABLE PaymentMethod (
    payment_method_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    fees_type ENUM('fixed', 'percentage') DEFAULT 'fixed',
    fees_amount DECIMAL(10, 2) DEFAULT 0,
    sort_order INT DEFAULT 0
);

CREATE TABLE Order_Table (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    subtotal DECIMAL(10, 2) NOT NULL,
    shipping_fee DECIMAL(10, 2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10, 2) NOT NULL,
    payment_method_id INT NOT NULL,
    payment_status ENUM('pending', 'paid', 'failed', 'refunded') DEFAULT 'pending',
    shipping_address_id INT NOT NULL,
    billing_address_id INT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES Customer(customer_id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES PaymentMethod(payment_method_id) ON DELETE CASCADE,
    FOREIGN KEY (shipping_address_id) REFERENCES CustomerAddress(address_id) ON DELETE CASCADE,
    FOREIGN KEY (billing_address_id) REFERENCES CustomerAddress(address_id) ON DELETE CASCADE
);

CREATE TABLE OrderItem (
    order_item_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    variant_id INT,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES Order_Table(order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES Product(product_id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES ProductVariant(variant_id) ON DELETE CASCADE
);

CREATE TABLE OrderStatusHistory (
    history_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled') NOT NULL,
    comment TEXT,
    staff_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES Order_Table(order_id) ON DELETE CASCADE,
    FOREIGN KEY (staff_id) REFERENCES Staff(staff_id) ON DELETE SET NULL
);

CREATE TABLE Shipment (
    shipment_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    tracking_number VARCHAR(100),
    carrier VARCHAR(100) NOT NULL,
    shipping_date TIMESTAMP,
    expected_delivery_date DATE,
    status ENUM('pending', 'shipped', 'delivered', 'failed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES Order_Table(order_id) ON DELETE CASCADE
);

-- 6. ระบบ Payment
CREATE TABLE Transaction (
    transaction_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    payment_method_id INT NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    amount DECIMAL(10, 2) NOT NULL,
    status ENUM('pending', 'completed', 'failed', 'refunded') DEFAULT 'pending',
    reference_number VARCHAR(100),
    gateway_response TEXT,
    gateway_transaction_id VARCHAR(100),
    FOREIGN KEY (order_id) REFERENCES Order_Table(order_id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES PaymentMethod(payment_method_id) ON DELETE CASCADE
);

CREATE TABLE Refund (
    refund_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    transaction_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    reason TEXT NOT NULL,
    status ENUM('pending', 'processed', 'rejected') DEFAULT 'pending',
    refund_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_by INT,
    notes TEXT,
    FOREIGN KEY (order_id) REFERENCES Order_Table(order_id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES Transaction(transaction_id) ON DELETE CASCADE,
    FOREIGN KEY (processed_by) REFERENCES Staff(staff_id) ON DELETE SET NULL
);

-- 7. ระบบออกเอกสารต่างๆ
CREATE TABLE Document (
    document_id INT AUTO_INCREMENT PRIMARY KEY,
    document_type ENUM('invoice', 'receipt', 'shipping_label', 'tax_invoice') NOT NULL,
    reference_id INT NOT NULL,
    reference_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(50) NOT NULL UNIQUE,
    generated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    file_path VARCHAR(255) NOT NULL,
    created_by INT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES Staff(staff_id) ON DELETE CASCADE
);

CREATE TABLE DocumentTemplate (
    template_id INT AUTO_INCREMENT PRIMARY KEY,
    document_type ENUM('invoice', 'receipt', 'shipping_label', 'tax_invoice') NOT NULL,
    name VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE
);