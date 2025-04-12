package product

// ProductRepository defines the interface for product operations
type ProductRepository interface {
	FindByID(id uint) (*Product, error)
	FindBySKU(sku string) (*Product, error)
	FindAll(page, limit int) ([]*Product, error)
	FindByCategory(categoryID uint, page, limit int) ([]*Product, error)
	FindWithFilters(filters map[string]interface{}, page, limit int) ([]*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id uint) error
	AddCategory(productID, categoryID uint) error
	RemoveCategory(productID, categoryID uint) error
	AddAttribute(productAttribute *ProductAttribute) error
	UpdateAttribute(productAttribute *ProductAttribute) error
	RemoveAttribute(productID, attributeID uint) error
}

// CategoryRepository defines the interface for category operations
type CategoryRepository interface {
	FindByID(id uint) (*Category, error)
	FindByName(name string) (*Category, error)
	FindAll() ([]*Category, error)
	FindAllWithProducts(page, limit int) ([]*Category, error)
	FindRootCategories() ([]*Category, error)
	GetCategoryTree() ([]*Category, error)
	FindChildren(categoryID uint) ([]*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint) error
}

// AttributeRepository defines the interface for attribute operations
type AttributeRepository interface {
	FindByID(id uint) (*Attribute, error)
	FindByName(name string) (*Attribute, error)
	FindAll() ([]*Attribute, error)
	FindByGroup(groupID uint) ([]*Attribute, error)
	Create(attribute *Attribute) error
	Update(attribute *Attribute) error
	Delete(id uint) error
	AddOption(attributeID uint, value string) (*AttributeOption, error)
	UpdateOption(option *AttributeOption) error
	RemoveOption(optionID uint) error
}

// AttributeGroupRepository defines the interface for attribute group operations
type AttributeGroupRepository interface {
	FindByID(id uint) (*AttributeGroup, error)
	FindByName(name string) (*AttributeGroup, error)
	FindAll() ([]*AttributeGroup, error)
	Create(group *AttributeGroup) error
	Update(group *AttributeGroup) error
	Delete(id uint) error
}

// ProductVariantRepository defines the interface for product variant operations
type ProductVariantRepository interface {
	FindByID(id uint) (*ProductVariant, error)
	FindBySKU(sku string) (*ProductVariant, error)
	FindByProduct(productID uint) ([]*ProductVariant, error)
	Create(variant *ProductVariant) error
	Update(variant *ProductVariant) error
	Delete(id uint) error
	AddAttribute(variantAttribute *VariantAttribute) error
	UpdateAttribute(variantAttribute *VariantAttribute) error
	RemoveAttribute(variantID, attributeID uint) error
}

// ProductImageRepository defines the interface for product image operations
type ProductImageRepository interface {
	FindByID(id uint) (*ProductImage, error)
	FindByProduct(productID uint) ([]*ProductImage, error)
	FindByVariant(variantID uint) ([]*ProductImage, error)
	Create(image *ProductImage) error
	Update(image *ProductImage) error
	Delete(id uint) error
	SetMainImage(imageID uint) error
}
