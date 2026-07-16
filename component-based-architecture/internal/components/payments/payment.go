package payments

type PaymentRequest struct {
	OrderID    string
	CustomerID string
	Amount     int
}

type CaptureResult struct{}

type RefundRequest struct {
	OrderID    string
	CustomerID string
	Amount     int
	Reason     string
}
