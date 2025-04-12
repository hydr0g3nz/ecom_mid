package order

import (
	"time"

	"github.com/hydr0g3nz/ecom_mid/domain/common"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeInvoice       DocumentType = "invoice"
	DocumentTypeReceipt       DocumentType = "receipt"
	DocumentTypeShippingLabel DocumentType = "shipping_label"
	DocumentTypeTaxInvoice    DocumentType = "tax_invoice"
)

// Document represents a generated document related to an order
type Document struct {
	common.Entity
	DocumentID     uint         `json:"document_id"`
	DocumentType   DocumentType `json:"document_type"`
	ReferenceID    uint         `json:"reference_id"`
	ReferenceType  string       `json:"reference_type"`
	DocumentNumber string       `json:"document_number"`
	GeneratedDate  time.Time    `json:"generated_date"`
	FilePath       string       `json:"file_path"`
	CreatedBy      uint         `json:"created_by"`
}

// DocumentTemplate represents a template for generating documents
type DocumentTemplate struct {
	common.Entity
	TemplateID   uint         `json:"template_id"`
	DocumentType DocumentType `json:"document_type"`
	Name         string       `json:"name"`
	Content      string       `json:"content"`
	IsDefault    bool         `json:"is_default"`
}
