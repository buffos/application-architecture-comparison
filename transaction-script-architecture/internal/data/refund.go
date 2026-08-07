package data

const (
	RefundStatusNotStarted = "NotStarted"
	RefundStatusPending    = "Pending"
	RefundStatusCompleted  = "Completed"
	RefundStatusFailed     = "Failed"
)

// Refund is a passive financial follow-up record for a return request.
type Refund struct {
	ID              string
	ReturnRequestID string
	OrderID         string
	Amount          int
	Status          string
	ProcessedBy     string
}
