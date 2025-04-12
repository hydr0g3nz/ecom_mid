package user

// StaffRepository defines the interface for staff operations
type StaffRepository interface {
	FindByID(id uint) (*Staff, error)
	FindByUsername(username string) (*Staff, error)
	FindByEmail(email string) (*Staff, error)
	FindAll(page, limit int) ([]*Staff, error)
	FindByRole(roleID uint) ([]*Staff, error)
	Create(staff *Staff) error
	Update(staff *Staff) error
	Delete(id uint) error
}

// CustomerRepository defines the interface for customer operations
type CustomerRepository interface {
	FindByID(id uint) (*Customer, error)
	FindByUsername(username string) (*Customer, error)
	FindByEmail(email string) (*Customer, error)
	FindAll(page, limit int) ([]*Customer, error)
	Create(customer *Customer) error
	Update(customer *Customer) error
	Delete(id uint) error
	AddAddress(customerID uint, address CustomerAddress) error
	UpdateAddress(address CustomerAddress) error
	DeleteAddress(addressID uint) error
	FindAddressesByCustomerID(customerID uint) ([]CustomerAddress, error)
}

// RoleRepository defines the interface for role operations
type RoleRepository interface {
	FindByID(id uint) (*Role, error)
	FindByName(name string) (*Role, error)
	FindAll() ([]*Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(id uint) error
	AddPermission(roleID uint, permissionID uint) error
	RemovePermission(roleID uint, permissionID uint) error
}

// PermissionRepository defines the interface for permission operations
type PermissionRepository interface {
	FindByID(id uint) (*Permission, error)
	FindByName(name string) (*Permission, error)
	FindAll() ([]*Permission, error)
	Create(permission *Permission) error
	Update(permission *Permission) error
	Delete(id uint) error
}
