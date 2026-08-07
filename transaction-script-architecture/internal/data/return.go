package data

import "time"

const (
	ReturnStatusRequested     = "Requested"
	ReturnStatusAccepted      = "Accepted"
	ReturnStatusRejected      = "Rejected"
	ReturnStatusRefundPending = "RefundPending"
	ReturnStatusRefunded      = "Refunded"
)

// ReturnLine is a passive request for a quantity from an order line.
type ReturnLine struct {
	OrderLineID     string
	SKU             string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

// ReturnRequest is a passive record coordinated by return transaction
// scripts. Actor fields are populated by a later auditability lesson.
type ReturnRequest struct {
	ID           string
	OrderID      string
	Status       string
	Reason       string
	Lines        []ReturnLine
	RefundID     string
	RefundStatus string
	RefundAmount int
	RequestedAt  time.Time
	RequestedBy  string
	ReviewedBy   string
	ProcessedBy  string
	ReviewNote   string
}
