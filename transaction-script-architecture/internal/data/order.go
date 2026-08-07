package data

const (
	OrderStatusPendingReservation  = "PendingReservation"
	OrderStatusBackordered         = "Backordered"
	OrderStatusReadyForPayment     = "ReadyForPayment"
	OrderStatusPaymentReview       = "PaymentReview"
	OrderStatusReadyForFulfillment = "ReadyForFulfillment"
	OrderStatusPartiallyShipped    = "PartiallyShipped"
	OrderStatusShipped             = "Shipped"
	OrderStatusCancelled           = "Cancelled"
	PaymentStatusNotRequired       = "NotRequired"
	PaymentStatusPending           = "Pending"
	PaymentStatusAccepted          = "Accepted"
	PaymentStatusFailed            = "Failed"
	PaymentStatusManualReview      = "ManualReview"
)

// OrderLine is a passive committed snapshot copied from a quote line.
type OrderLine struct {
	ID                  string
	SKU                 string
	ProductNameSnapshot string
	ProductCategory     string
	OrderedQuantity     int
	ReservedQuantity    int
	ShippedQuantity     int
	ReturnedQuantity    int
	UnitPrice           int
	DiscountAmount      int
	LineTotal           int
}

// Order is a passive record coordinated by transaction scripts.
type Order struct {
	ID                 string
	SourceQuoteID      string
	CustomerID         string
	Status             string
	RequestedBy        string
	Lines              []OrderLine
	Total              int
	PaymentID          string
	PaymentStatus      string
	CancelledBy        string
	CancellationReason string
}
