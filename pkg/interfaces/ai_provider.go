package interfaces

// AIProvider defines the interface for AI service providers
type AIProvider interface {
	// GetReceiptInvoiceInfo extracts structured information from receipt/invoice text
	GetReceiptInvoiceInfo(content string) (*ReceiptInvoiceInfo, error)
}

// IdField represents an identification field found in the document
type IdField struct {
	// Name is the type/name of the identifier (e.g., "Invoice Number", "Receipt Number", "Customer ID")
	Name string `json:"name" jsonschema_description:"The type or name of the identifier (e.g., 'Invoice Number', 'Receipt Number', 'Customer ID')"`
	
	// Value is the actual identifier value
	Value string `json:"value" jsonschema_description:"The actual identifier value"`
}

// ReceiptInvoiceInfo represents the structured information extracted from a receipt or invoice
type ReceiptInvoiceInfo struct {
	// DocumentType is always present and classifies the document
	DocumentType string `json:"document_type" jsonschema:"enum=None,enum=Invoice,enum=Receipt" jsonschema_description:"Classification of the document as None, Invoice, or Receipt"`
	
	// Description is a mandatory accountant-friendly categorization of the document/service
	Description string `json:"description" jsonschema:"maxLength=50" jsonschema_description:"Mandatory accountant-friendly description: for None documents describe what it's about, for Invoice/Receipt provide generic service category (e.g., 'AI Services', 'Cloud Services'). Max 50 characters."`
	
	// Company is the entity offering the service and requesting payment (nullable)
	Company *string `json:"company" jsonschema_description:"The company that owns the service being offered and is requesting payment, null if not found"`
	
	// DateIssued is the date the invoice/receipt was issued (nullable)
	DateIssued *string `json:"date_issued" jsonschema_description:"The date the document was issued in YYYY-MM-DD format, null if not found"`
	
	// ServiceDescription is a description of the service or items paid for (nullable)
	ServiceDescription *string `json:"service_description" jsonschema_description:"Description of the service or items paid for, null if not found"`
	
	// SECentAmount is the amount in Swedish cents (nullable)
	// Last 2 digits are öre (cents), rest is kronor
	// For example: 9537 = 95.37 SEK
	SECentAmount *int `json:"se_cent_amount" jsonschema_description:"Amount in Swedish cents (öre), where last 2 digits are cents and rest is kronor, null if not found"`
	
	// OriginalAmount is the total amount in the original currency (nullable)
	OriginalAmount *float64 `json:"original_amount" jsonschema_description:"The total amount in the original currency, null if not found"`
	
	// OriginalCurrency is the ISO 3-letter currency code (nullable)
	OriginalCurrency *string `json:"original_currency" jsonschema_description:"The ISO 3-letter currency code (e.g., 'EUR', 'USD', 'SEK'), null if not found"`
	
	// OriginalVatAmount is the VAT amount in the original currency (nullable)
	OriginalVatAmount *float64 `json:"original_vat_amount" jsonschema_description:"The VAT/tax amount in the original currency, null if not found"`
	
	// IdFields is a list of identification fields found in the document
	IdFields []IdField `json:"id_fields" jsonschema_description:"List of identification fields found in the document (invoice numbers, receipt numbers, customer IDs, etc.). Can be empty."`
	
	// SuggestedFileName is a generated filename based on extracted data (populated post-processing)
	// Format: <date>-<company>-<description>-<amount>SEK (lowercase, non-alphanumeric chars become _)
	SuggestedFileName string `json:"suggested_filename" jsonschema:"-"`
}