package product

import (
	"github.com/hydr0g3nz/ecom_mid/domain/common"
	"github.com/hydr0g3nz/ecom_mid/domain/common/vo"
)

// CategoryStatus represents the status of a category
type CategoryStatus string

const (
	CategoryStatusActive   CategoryStatus = "active"
	CategoryStatusInactive CategoryStatus = "inactive"
)

// Category represents a product category
type Category struct {
	common.Entity
	CategoryID       uint           `json:"category_id"`
	ParentCategoryID *uint          `json:"parent_category_id,omitempty"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	Image            string         `json:"image"`
	Status           CategoryStatus `json:"status"`
	DisplayOrder     int            `json:"display_order"`
	Children         []Category     `json:"children,omitempty"`
	Parent           *Category      `json:"parent,omitempty"`
}

// IsActive checks if the category is active
func (c *Category) IsActive() bool {
	return c.Status == CategoryStatusActive
}

// Activate sets the category status to active
func (c *Category) Activate() {
	c.Status = CategoryStatusActive
}

// Deactivate sets the category status to inactive
func (c *Category) Deactivate() {
	c.Status = CategoryStatusInactive
}

// IsRoot checks if the category is a root category (no parent)
func (c *Category) IsRoot() bool {
	return c.ParentCategoryID == nil
}

// IsLeaf checks if the category is a leaf category (no children)
func (c *Category) IsLeaf() bool {
	return len(c.Children) == 0
}

// AddChild adds a child category
func (c *Category) AddChild(child Category) {
	// Set this category as parent
	child.ParentCategoryID = &c.CategoryID
	child.Parent = c
	
	c.Children = append(c.Children, child)
}
