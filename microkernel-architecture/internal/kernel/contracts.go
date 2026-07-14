package kernel

import "errors"

var ErrPluginAlreadyRegistered = errors.New("plugin already registered")
var ErrCustomerDirectoryNotRegistered = errors.New("customer directory capability not registered")
var ErrProductCatalogNotRegistered = errors.New("product catalog capability not registered")
var ErrApprovalPolicyNotRegistered = errors.New("approval policy capability not registered")
var ErrQuoteServiceNotRegistered = errors.New("quote service capability not registered")
var ErrQuoteReaderNotRegistered = errors.New("quote reader capability not registered")
var ErrApprovedQuoteProviderNotRegistered = errors.New("approved quote provider capability not registered")
var ErrInventoryReservationNotRegistered = errors.New("inventory reservation capability not registered")
var ErrOrderServiceNotRegistered = errors.New("order service capability not registered")

type Plugin interface {
	ID() string
	Register(host *Host) error
}

type CustomerDirectory interface {
	RequireActiveCustomer(id string) error
}

type Product struct {
	SKU       string
	Name      string
	Category  string
	UnitPrice int
}

type ProductCatalog interface {
	GetProductForQuote(sku string) (Product, error)
}

type QuoteSubmissionLine struct {
	ProductCategory string
}

type QuoteSubmission struct {
	Lines []QuoteSubmissionLine
}

type ApprovalPolicy interface {
	RequiresApproval(submission QuoteSubmission) bool
}

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type AddQuoteLineCommand struct {
	QuoteID    string
	ProductSKU string
	Quantity   int
}

type AddQuoteLineResult struct {
	QuoteID    string
	LineCount  int
	TotalItems int
	Status     string
}

type SubmitQuoteCommand struct {
	QuoteID string
}

type SubmitQuoteResult struct {
	QuoteID    string
	LineCount  int
	TotalItems int
	Status     string
}

type ApproveQuoteCommand struct {
	QuoteID string
}

type ApproveQuoteResult struct {
	QuoteID    string
	LineCount  int
	TotalItems int
	Status     string
}

type QuoteService interface {
	CreateDraftQuote(command CreateDraftQuoteCommand) (CreateDraftQuoteResult, error)
	AddQuoteLine(command AddQuoteLineCommand) (AddQuoteLineResult, error)
	SubmitQuote(command SubmitQuoteCommand) (SubmitQuoteResult, error)
	ApproveQuote(command ApproveQuoteCommand) (ApproveQuoteResult, error)
}

type GetQuoteQuery struct {
	QuoteID string
}

type QuoteDetails struct {
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
	TotalItems int
}

type QuoteReader interface {
	GetQuote(query GetQuoteQuery) (QuoteDetails, error)
}

type ApprovedQuote struct {
	QuoteID    string
	CustomerID string
	Lines      []ApprovedQuoteLine
}

type ApprovedQuoteLine struct {
	ProductSKU      string
	ProductName     string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

type ApprovedQuoteProvider interface {
	GetApprovedQuoteForOrder(quoteID string) (ApprovedQuote, error)
}

type InventoryReservationItem struct {
	ProductSKU string
	Quantity   int
}

type InventoryReservation interface {
	Reserve(items []InventoryReservationItem) error
}

type ConvertQuoteToOrderCommand struct {
	QuoteID string
}

type ConvertQuoteToOrderResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type OrderService interface {
	ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error)
}
