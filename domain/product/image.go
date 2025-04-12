package product

import (
	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// ProductImage represents an image associated with a product or variant
type ProductImage struct {
	common.Entity
	ImageID   uint   `json:"image_id"`
	ProductID uint   `json:"product_id"`
	VariantID *uint  `json:"variant_id,omitempty"`
	ImagePath string `json:"image_path"`
	AltText   string `json:"alt_text"`
	SortOrder int    `json:"sort_order"`
	IsMain    bool   `json:"is_main"`
}
