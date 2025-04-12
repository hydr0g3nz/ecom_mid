package product

import (
	"errors"
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusDraft    ProductStatus = "draft"
)

// Product represents a product in the system
type Product struct {
	common.Entity
	ProductID          uint           `json:"product_id"`
	SKU                string         `json:"sku"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	ShortDescription   string         `json:"short_description"`
	Price              vo.Money       `json:"price"`
	SpecialPrice       *vo.Money      `json:"special_price,omitempty"`
	SpecialPriceStart  *time.Time     `json:"special_price_start,omitempty"`
	SpecialPriceEnd    *time.Time     `json:"special_price_end,omitempty"`
	Cost               *vo.Money      `json:"cost,omitempty"`
	Weight             float64        `json:"weight"`
	Length             float64        `json:"length"`
	Width              float64        `json:"width"`
	Height             float64        `json:"height"`
	Status             ProductStatus  `json:"status"`
	MetaTitle          string         `json:"meta_title"`
	MetaDescription    string         `json:"meta_description"`
	MetaKeywords       string         `json:"meta_keywords"`
	Categories         []Category     `json:"categories"`
	Attributes         []ProductAttribute `json:"attributes"`
	Variants           []ProductVariant   `json:"variants"`
	Images             []ProductImage     `json:"images"`
}

// NewProduct creates a new product with validation
func NewProduct(sku, name, description string, price vo.Money) (*Product, error) {
	if sku == "" {
		return nil, errors.New("sku cannot be empty")
	}
	
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	
	if price.IsZero() || price.IsNegative() {
		return nil, errors.New("price must be greater than zero")
	}
	
	return &Product{
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		Status:      ProductStatusDraft,
		Categories:  []Category{},
		Attributes:  []ProductAttribute{},
		Variants:    []ProductVariant{},
		Images:      []ProductImage{},
	}, nil
}

// IsActive checks if the product is active
func (p *Product) IsActive() bool {
	return p.Status == ProductStatusActive
}

// Activate sets the product status to active
func (p *Product) Activate() {
	p.Status = ProductStatusActive
}

// Deactivate sets the product status to inactive
func (p *Product) Deactivate() {
	p.Status = ProductStatusInactive
}

// GetCurrentPrice returns the current applicable price (special or regular)
func (p *Product) GetCurrentPrice() vo.Money {
	if p.HasActiveSpecialPrice() {
		return *p.SpecialPrice
	}
	return p.Price
}

// HasActiveSpecialPrice checks if the product has an active special price
func (p *Product) HasActiveSpecialPrice() bool {
	if p.SpecialPrice == nil || p.SpecialPrice.IsZero() {
		return false
	}
	
	now := time.Now()
	
	// If no dates specified, special price is always active
	if p.SpecialPriceStart == nil && p.SpecialPriceEnd == nil {
		return true
	}
	
	// Check start date if specified
	if p.SpecialPriceStart != nil && now.Before(*p.SpecialPriceStart) {
		return false
	}
	
	// Check end date if specified
	if p.SpecialPriceEnd != nil && now.After(*p.SpecialPriceEnd) {
		return false
	}
	
	return true
}

// SetSpecialPrice sets a special price with optional date range
func (p *Product) SetSpecialPrice(price vo.Money, startDate, endDate *time.Time) error {
	if price.IsNegative() {
		return errors.New("special price cannot be negative")
	}
	
	// Validate date range if both are provided
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return errors.New("end date cannot be before start date")
	}
	
	p.SpecialPrice = &price
	p.SpecialPriceStart = startDate
	p.SpecialPriceEnd = endDate
	return nil
}

// AddCategory adds a category to the product
func (p *Product) AddCategory(category Category) {
	// Check if category already exists
	for _, c := range p.Categories {
		if c.CategoryID == category.CategoryID {
			return
		}
	}
	p.Categories = append(p.Categories, category)
}

// RemoveCategory removes a category from the product
func (p *Product) RemoveCategory(categoryID uint) {
	updatedCategories := []Category{}
	for _, category := range p.Categories {
		if category.CategoryID != categoryID {
			updatedCategories = append(updatedCategories, category)
		}
	}
	p.Categories = updatedCategories
}

// AddAttribute adds or updates a product attribute
func (p *Product) AddAttribute(attr ProductAttribute) {
	// Update existing attribute if found
	for i, a := range p.Attributes {
		if a.AttributeID == attr.AttributeID {
			p.Attributes[i] = attr
			return
		}
	}
	
	// Add new attribute
	p.Attributes = append(p.Attributes, attr)
}

// RemoveAttribute removes an attribute from the product
func (p *Product) RemoveAttribute(attributeID uint) {
	updatedAttributes := []ProductAttribute{}
	for _, attr := range p.Attributes {
		if attr.AttributeID != attributeID {
			updatedAttributes = append(updatedAttributes, attr)
		}
	}
	p.Attributes = updatedAttributes
}

// AddVariant adds a variant to the product
func (p *Product) AddVariant(variant ProductVariant) error {
	// Check if SKU is unique
	for _, v := range p.Variants {
		if v.SKU == variant.SKU {
			return errors.New("variant SKU must be unique")
		}
	}
	
	p.Variants = append(p.Variants, variant)
	return nil
}

// RemoveVariant removes a variant from the product
func (p *Product) RemoveVariant(variantID uint) {
	updatedVariants := []ProductVariant{}
	for _, variant := range p.Variants {
		if variant.VariantID != variantID {
			updatedVariants = append(updatedVariants, variant)
		}
	}
	p.Variants = updatedVariants
}

// AddImage adds an image to the product
func (p *Product) AddImage(image ProductImage) {
	// If main image is being added, set existing main image to non-main
	if image.IsMain {
		for i := range p.Images {
			p.Images[i].IsMain = false
		}
	}
	
	// If this is the first image, make it the main image
	if len(p.Images) == 0 {
		image.IsMain = true
	}
	
	p.Images = append(p.Images, image)
}

// RemoveImage removes an image from the product
func (p *Product) RemoveImage(imageID uint) error {
	removedMainImage := false
	updatedImages := []ProductImage{}
	
	for _, img := range p.Images {
		if img.ImageID != imageID {
			updatedImages = append(updatedImages, img)
		} else if img.IsMain {
			removedMainImage = true
		}
	}
	
	p.Images = updatedImages
	
	// If we removed the main image and have other images, set the first one as main
	if removedMainImage && len(p.Images) > 0 {
		p.Images[0].IsMain = true
	}
	
	return nil
}

// GetMainImage returns the main product image or nil if none exists
func (p *Product) GetMainImage() *ProductImage {
	for _, img := range p.Images {
		if img.IsMain {
			return &img
		}
	}
	
	// If no main image is set but we have images, return the first one
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	
	return nil
}
