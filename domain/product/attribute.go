package product

import (
	"errors"
	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// AttributeType represents the type of attribute
type AttributeType string

const (
	AttributeTypeText    AttributeType = "text"
	AttributeTypeSelect  AttributeType = "select"
	AttributeTypeBoolean AttributeType = "boolean"
	AttributeTypeNumber  AttributeType = "number"
	AttributeTypeDate    AttributeType = "date"
)

// AttributeGroup represents a group of related attributes
type AttributeGroup struct {
	common.Entity
	AttributeGroupID uint        `json:"attribute_group_id"`
	Name             string      `json:"name"`
	Attributes       []Attribute `json:"attributes,omitempty"`
}

// Attribute represents a product attribute
type Attribute struct {
	common.Entity
	AttributeID      uint           `json:"attribute_id"`
	AttributeGroupID uint           `json:"attribute_group_id"`
	Name             string         `json:"name"`
	Type             AttributeType  `json:"type"`
	IsFilterable     bool           `json:"is_filterable"`
	IsRequired       bool           `json:"is_required"`
	Options          []AttributeOption `json:"options,omitempty"`
	Group            *AttributeGroup  `json:"group,omitempty"`
}

// AttributeOption represents a predefined option for select-type attributes
type AttributeOption struct {
	common.Entity
	OptionID    uint   `json:"option_id"`
	AttributeID uint   `json:"attribute_id"`
	Value       string `json:"value"`
}

// ProductAttribute represents an attribute value for a specific product
type ProductAttribute struct {
	common.Entity
	ProductAttributeID uint           `json:"product_attribute_id"`
	ProductID          uint           `json:"product_id"`
	AttributeID        uint           `json:"attribute_id"`
	Value              string         `json:"value,omitempty"`
	OptionID           *uint          `json:"option_id,omitempty"`
	Attribute          *Attribute     `json:"attribute,omitempty"`
	Option             *AttributeOption `json:"option,omitempty"`
}

// AddOption adds an option to an attribute
func (a *Attribute) AddOption(value string) (AttributeOption, error) {
	if a.Type != AttributeTypeSelect {
		return AttributeOption{}, errors.New("options can only be added to select attributes")
	}
	
	// Check for duplicate option values
	for _, option := range a.Options {
		if option.Value == value {
			return option, errors.New("option value already exists")
		}
	}
	
	option := AttributeOption{
		AttributeID: a.AttributeID,
		Value: value,
	}
	
	a.Options = append(a.Options, option)
	return option, nil
}

// RemoveOption removes an option from an attribute
func (a *Attribute) RemoveOption(optionID uint) {
	updatedOptions := []AttributeOption{}
	for _, option := range a.Options {
		if option.OptionID != optionID {
			updatedOptions = append(updatedOptions, option)
		}
	}
	a.Options = updatedOptions
}

// ValidateAttributeValue validates a value against the attribute's type
func (a *Attribute) ValidateAttributeValue(value string, optionID *uint) error {
	// For select attributes, an option ID must be provided
	if a.Type == AttributeTypeSelect {
		if optionID == nil {
			return errors.New("option ID is required for select attributes")
		}
		
		// Verify the option ID belongs to this attribute
		valid := false
		for _, option := range a.Options {
			if option.OptionID == *optionID {
				valid = true
				break
			}
		}
		
		if !valid {
			return errors.New("invalid option ID for this attribute")
		}
		
		return nil
	}
	
	// For non-select attributes, no option ID should be provided
	if optionID != nil {
		return errors.New("option ID should only be provided for select attributes")
	}
	
	// For boolean attributes, value should be "true" or "false"
	if a.Type == AttributeTypeBoolean && value != "true" && value != "false" {
		return errors.New("boolean attribute value must be 'true' or 'false'")
	}
	
	// Other validation could be added for number and date types
	
	return nil
}
