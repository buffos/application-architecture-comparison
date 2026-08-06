package payments

type PaymentRequest struct {
	OrderID    string
	CustomerID string
	Amount     int
}

type CaptureOutcome string

const (
	CaptureOutcomeApproved CaptureOutcome = "Approved"
	CaptureOutcomeReview   CaptureOutcome = "Review"
)

type CaptureResult struct {
	Outcome CaptureOutcome
}

type RefundRequest struct {
	OrderID    string
	CustomerID string
	Amount     int
	Reason     string
}
