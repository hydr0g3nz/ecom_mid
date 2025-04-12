package product

import (
	"errors"
	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
)

// ProductVariantStatus represents the status of a product variant
type ProductVariantStatus string

const (
	VariantStatusActive   ProductVariantStatus = "active"
	VariantStatusInactive ProductVariantStatus = "inactive"
)

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	common.Entity
	VariantID     uint                 `json:"variant_id"`
	ProductID     uint                 `json:"product_id"`
	SKU           string               `json:"sku"`
	Price         vo.Money             `json:"price"`
	SpecialPrice  *vo.Money            `json:"special_price,omitempty"`
	StockQuantity int                  `json:"stock_quantity"`
	Weight        float64              `json:"weight"`
	Status        ProductVariantStatus `json:"status"`
	Attributes    []VariantAttribute   `json:"attributes,omitempty"`
	Images        []ProductImage       `json:"images,omitempty"`
}

// VariantAttribute represents an attribute that defines a variant
type VariantAttribute struct {
	common.Entity
	VariantAttributeID uint            `json:"variant_attribute_id"`
	VariantID          uint            `json:"variant_id"`
	AttributeID        uint            `json:"attribute_id"`
	OptionID           uint            `json:"option_id"`
	Attribute          *Attribute      `json:"attribute,omitempty"`
	Option             *AttributeOption `json:"option,omitempty"`
}

// IsActive checks if the variant is active
func (v *ProductVariant) IsActive() bool {
	return v.Status == VariantStatusActive
}

// Activate sets the variant status to active
func (v *ProductVariant) Activate() {
	v.Status = VariantStatusActive
}

// Deactivate sets the variant status to inactive
func (v *ProductVariant) Deactivate() {
	v.Status = VariantStatusInactive
}

// IsInStock checks if the variant has stock available
func (v *ProductVariant) IsInStock() bool {
	return v.StockQuantity > 0
}

// AddStock adds stock quantity to the variant
func (v *ProductVariant) AddStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	v.StockQuantity += quantity
	return nil
}

// RemoveStock removes stock quantity from the variant
func (v *ProductVariant) RemoveStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	
	if v.StockQuantity < quantity {
		return errors.New("insufficient stock")
	}
	
	v.StockQuantity -= quantity
	return nil
}

// GetCurrentPrice returns the current applicable price (special or regular)
func (v *ProductVariant) GetCurrentPrice() vo.Money {
	if v.SpecialPrice != nil && !v.SpecialPrice.IsZero() {
		return *v.SpecialPrice
	}
	return v.Price
}

// AddAttribute adds a variant attribute
func (v *ProductVariant) AddAttribute(attr VariantAttribute) error {
	// Check if attribute already exists for this variant
	for _, a := range v.Attributes {
		if a.AttributeID == attr.AttributeID {
			return errors.New("variant already has this attribute")
		}
	}
	
	v.Attributes = append(v.Attributes, attr)
	return nil
}
