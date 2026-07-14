package kernel

import (
	"errors"
	"time"
)

var ErrPluginAlreadyRegistered = errors.New("plugin already registered")
var ErrCustomerDirectoryNotRegistered = errors.New("customer directory capability not registered")
var ErrProductCatalogNotRegistered = errors.New("product catalog capability not registered")
var ErrApprovalPolicyNotRegistered = errors.New("approval policy capability not registered")
var ErrQuoteServiceNotRegistered = errors.New("quote service capability not registered")
var ErrQuoteReaderNotRegistered = errors.New("quote reader capability not registered")
var ErrApprovedQuoteProviderNotRegistered = errors.New("approved quote provider capability not registered")
var ErrInventoryReservationNotRegistered = errors.New("inventory reservation capability not registered")
var ErrInventoryReleaseNotRegistered = errors.New("inventory release capability not registered")
var ErrInventoryRestockNotRegistered = errors.New("inventory restock capability not registered")
var ErrPaymentCaptureNotRegistered = errors.New("payment capture capability not registered")
var ErrPaymentRefundNotRegistered = errors.New("payment refund capability not registered")
var ErrShipmentCreationNotRegistered = errors.New("shipment creation capability not registered")
var ErrOrderServiceNotRegistered = errors.New("order service capability not registered")
var ErrOrderReaderNotRegistered = errors.New("order reader capability not registered")
var ErrReturnableOrderProviderNotRegistered = errors.New("returnable order provider capability not registered")
var ErrReturnEligibilityPolicyNotRegistered = errors.New("return eligibility policy capability not registered")
var ErrReturnServiceNotRegistered = errors.New("return service capability not registered")
var ErrReturnReaderNotRegistered = errors.New("return reader capability not registered")
var ErrClockNotRegistered = errors.New("clock capability not registered")
var ErrIdempotencyStoreNotRegistered = errors.New("idempotency store capability not registered")
var ErrIdempotencyKeyRequired = errors.New("idempotency key is required")

type Plugin interface {
	ID() string
	Register(host *Host) error
}

type CustomerDirectory interface {
	RequireActiveCustomer(id string) error
}

type Product struct {
	SKU              string
	Name             string
	Category         string
	UnitPrice        int
	ReturnWindowDays int
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
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
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

type InventoryRelease interface {
	Release(items []InventoryReservationItem) error
}

type InventoryRestock interface {
	Restock(items []InventoryReservationItem) error
}

type PaymentCapture interface {
	Capture(orderID string, amount int) error
}

type PaymentRefund interface {
	Refund(orderID string, amount int) error
}

type Clock interface {
	Now() time.Time
}

type IdempotencyResult struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type IdempotencyStore interface {
	Find(key string) (IdempotencyResult, bool, error)
	Save(key string, result IdempotencyResult) error
}

type ReturnEligibilityReview struct {
	Reason      string
	ShippedAt   time.Time
	RequestedAt time.Time
	Lines       []ReturnEligibilityLine
}

type ReturnEligibilityLine struct {
	ReturnWindowDays int
}

type ReturnEligibilityPolicy interface {
	Allows(review ReturnEligibilityReview) bool
}

type ShipmentLine struct {
	ProductSKU string
	Quantity   int
}

type CreateShipmentRecord struct {
	OrderID    string
	CustomerID string
	Lines      []ShipmentLine
}

type ShipmentCreationResult struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	LineCount  int
}

type ShipmentCreation interface {
	CreateShipment(record CreateShipmentRecord) (ShipmentCreationResult, error)
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

type CapturePaymentCommand struct {
	OrderID string
}

type CapturePaymentResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CreateShipmentCommand struct {
	OrderID string
}

type CreateShipmentResult struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CancelOrderCommand struct {
	OrderID string
}

type CancelOrderResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type GetOrderQuery struct {
	OrderID string
}

type OrderDetails struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ListOrdersQuery struct {
	Status string
}

type OrderSummary struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ReturnableOrder struct {
	OrderID    string
	CustomerID string
	ShippedAt  time.Time
	Lines      []ReturnableOrderLine
}

type ReturnableOrderLine struct {
	ProductSKU       string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

type ReturnableOrderProvider interface {
	GetReturnableOrder(orderID string) (ReturnableOrder, error)
}

type RequestReturnCommand struct {
	OrderID     string
	Reason      string
	RequestedBy string
}

type RequestReturnResult struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type AcceptReturnCommand struct {
	ReturnRequestID string
	IdempotencyKey  string
	ReviewedBy      string
	ProcessedBy     string
	ReviewNote      string
}

type AcceptReturnResult struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type RejectReturnCommand struct {
	ReturnRequestID string
	IdempotencyKey  string
	ReviewedBy      string
	ReviewNote      string
}

type RejectReturnResult struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type GetReturnRequestQuery struct {
	ReturnRequestID string
}

type ReturnRequestDetails struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	Reason          string
	LineCount       int
	RequestedBy     string
	ReviewedBy      string
	ProcessedBy     string
	ReviewNote      string
}

type ListReturnRequestsQuery struct {
	Status string
}

type ReturnRequestSummary struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type OrderService interface {
	ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error)
	CapturePayment(command CapturePaymentCommand) (CapturePaymentResult, error)
	CreateShipment(command CreateShipmentCommand) (CreateShipmentResult, error)
	CancelOrder(command CancelOrderCommand) (CancelOrderResult, error)
}

type OrderReader interface {
	GetOrder(query GetOrderQuery) (OrderDetails, error)
	ListOrders(query ListOrdersQuery) ([]OrderSummary, error)
}

type ReturnService interface {
	RequestReturn(command RequestReturnCommand) (RequestReturnResult, error)
	AcceptReturn(command AcceptReturnCommand) (AcceptReturnResult, error)
	RejectReturn(command RejectReturnCommand) (RejectReturnResult, error)
}

type ReturnReader interface {
	GetReturnRequest(query GetReturnRequestQuery) (ReturnRequestDetails, error)
	ListReturnRequests(query ListReturnRequestsQuery) ([]ReturnRequestSummary, error)
}
