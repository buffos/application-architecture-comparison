package payments

type PaymentRequest struct {
	OrderID    string
	CustomerID string
	Amount     int
}

const (
	CaptureOutcomeApproved = "Approved"
	CaptureOutcomeReview   = "Review"
)

type CaptureResult struct {
	Outcome string
}

type RefundRequest struct {
	OrderID    string
	CustomerID string
	Amount     int
	Reason     string
}
