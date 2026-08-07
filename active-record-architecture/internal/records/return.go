package records

import (
	"errors"
	"time"
)

const (
	DefaultReturnWindowDays = 30
	ReturnStatusRequested   = "Requested"
	ReturnStatusAccepted    = "Accepted"
	ReturnStatusRejected    = "Rejected"
	ReturnStatusRefunded    = "Refunded"
	RefundStatusNotStarted  = "NotStarted"
	RefundStatusCompleted   = "Completed"
)

var (
	ErrOrderNotReturnable     = errors.New("order is not returnable")
	ErrReturnIDRequired       = errors.New("return id is required")
	ErrReturnNotFound         = errors.New("return request not found")
	ErrReturnLinesInvalid     = errors.New("return lines are invalid")
	ErrReturnNotAcceptable    = errors.New("return is not requested")
	ErrReturnNotRejectable    = errors.New("return is not requested")
	ErrReturnNotRefundable    = errors.New("return refund cannot be completed")
	ErrReturnNotEligible      = errors.New("return is not eligible")
	ErrReturnOrderMissing     = errors.New("return order not found")
	ErrReturnStockMissing     = errors.New("return stock record not found")
	ErrActorRequired          = errors.New("return actor is required")
	ErrIdempotencyKeyRequired = errors.New("idempotency key is required")
)

// ReturnLine is a passive request for a quantity from an order-line snapshot.
type ReturnLine struct {
	OrderLineID      string
	SKU              string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

// ReturnRequest is an Active Record for the requested reverse flow. Its
// business transitions are added in later lessons.
type ReturnRequest struct {
	db *Database

	ID           string
	OrderID      string
	Status       string
	Reason       string
	ReviewNote   string
	RequestedBy  string
	ReviewedBy   string
	ProcessedBy  string
	Lines        []ReturnLine
	RefundID     string
	RefundStatus string
	RefundAmount int
	RequestedAt  time.Time
}

// Accept records a positive review decision. Side effects wait for
// CompleteRefund so the review boundary is explicit.
func (request *ReturnRequest) Accept(reviewedBy string, idempotencyKey string) (*ReturnRequest, error) {
	return request.AcceptAt(time.Now(), reviewedBy, idempotencyKey)
}

// AcceptAt is the deterministic form of Accept used by tests and
// demonstrations.
func (request *ReturnRequest) AcceptAt(now time.Time, reviewedBy string, idempotencyKey string) (*ReturnRequest, error) {
	if request == nil || request.db == nil {
		return nil, ErrDatabaseRequired
	}
	if reviewedBy == "" {
		return nil, ErrActorRequired
	}
	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired
	}
	if existing, found := findIdempotentReturn(request.db, "accept-return", idempotencyKey); found {
		return existing, nil
	}
	if request.Status != ReturnStatusRequested {
		return nil, ErrReturnNotAcceptable
	}
	order, err := FindOrder(request.db, request.OrderID)
	if err != nil {
		return nil, ErrReturnOrderMissing
	}
	if decision := request.EvaluateEligibilityAt(order, now); !decision.Eligible {
		return nil, ErrReturnNotEligible
	}
	request.Status = ReturnStatusAccepted
	request.ReviewedBy = reviewedBy
	if err := request.Save(); err != nil {
		return nil, err
	}
	saveIdempotentReturn(request.db, "accept-return", idempotencyKey, request)
	return request, nil
}

// Reject records a negative review decision without changing the order,
// stock, or refund records.
func (request *ReturnRequest) Reject(reviewedBy string, reviewNote string, idempotencyKey string) (*ReturnRequest, error) {
	if request == nil || request.db == nil {
		return nil, ErrDatabaseRequired
	}
	if reviewedBy == "" {
		return nil, ErrActorRequired
	}
	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired
	}
	if existing, found := findIdempotentReturn(request.db, "reject-return", idempotencyKey); found {
		return existing, nil
	}
	if request.Status != ReturnStatusRequested {
		return nil, ErrReturnNotRejectable
	}
	request.Status = ReturnStatusRejected
	request.ReviewedBy = reviewedBy
	request.ReviewNote = reviewNote
	if err := request.Save(); err != nil {
		return nil, err
	}
	saveIdempotentReturn(request.db, "reject-return", idempotencyKey, request)
	return request, nil
}

// CompleteRefund applies the accepted return's reverse side effects and
// completes its linked refund.
func (request *ReturnRequest) CompleteRefund(processedBy string, idempotencyKey string) (*ReturnRequest, error) {
	if request == nil || request.db == nil {
		return nil, ErrDatabaseRequired
	}
	if processedBy == "" {
		return nil, ErrActorRequired
	}
	if idempotencyKey == "" {
		return nil, ErrIdempotencyKeyRequired
	}
	if existing, found := findIdempotentReturn(request.db, "complete-refund", idempotencyKey); found {
		return existing, nil
	}
	if request.Status != ReturnStatusAccepted {
		return nil, ErrReturnNotRefundable
	}

	order, err := FindOrder(request.db, request.OrderID)
	if err != nil {
		return nil, ErrReturnOrderMissing
	}
	refund, err := FindRefund(request.db, request.RefundID)
	if err != nil {
		return nil, ErrReturnNotRefundable
	}
	if refund.Status != RefundStatusNotStarted {
		return nil, ErrReturnNotRefundable
	}

	type validatedLine struct {
		orderLineIndex int
		stock          *StockRecord
		quantity       int
	}
	validated := make([]validatedLine, 0, len(request.Lines))
	for _, returnLine := range request.Lines {
		if returnLine.Quantity <= 0 {
			return nil, ErrReturnLinesInvalid
		}

		orderLineIndex := -1
		for index, orderLine := range order.Lines {
			if orderLine.ID != returnLine.OrderLineID {
				continue
			}
			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if returnLine.Quantity > remaining {
				return nil, ErrReturnLinesInvalid
			}
			orderLineIndex = index
			break
		}
		if orderLineIndex < 0 {
			return nil, ErrReturnLinesInvalid
		}

		stock, err := FindStock(request.db, order.Lines[orderLineIndex].SKU)
		if err != nil {
			return nil, ErrReturnStockMissing
		}
		validated = append(validated, validatedLine{
			orderLineIndex: orderLineIndex,
			stock:          stock,
			quantity:       returnLine.Quantity,
		})
	}

	for _, line := range validated {
		order.Lines[line.orderLineIndex].ReturnedQuantity += line.quantity
		line.stock.OnHand += line.quantity
		if err := line.stock.Save(); err != nil {
			return nil, err
		}
	}

	refund.Status = RefundStatusCompleted
	refund.ProcessedBy = processedBy
	if err := refund.Save(); err != nil {
		return nil, err
	}
	request.Status = ReturnStatusRefunded
	request.RefundStatus = RefundStatusCompleted
	request.ProcessedBy = processedBy
	if err := order.Save(); err != nil {
		return nil, err
	}
	if err := request.Save(); err != nil {
		return nil, err
	}
	saveIdempotentReturn(request.db, "complete-refund", idempotencyKey, request)
	return request, nil
}

// FindReturnRequest loads a return request Active Record from the returns
// table.
func FindReturnRequest(db *Database, id string) (*ReturnRequest, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrReturnIDRequired
	}

	row, ok := db.returns[id]
	if !ok {
		return nil, ErrReturnNotFound
	}
	return &ReturnRequest{
		db:           db,
		ID:           row.ID,
		OrderID:      row.OrderID,
		Status:       row.Status,
		Reason:       row.Reason,
		ReviewNote:   row.ReviewNote,
		RequestedBy:  row.RequestedBy,
		ReviewedBy:   row.ReviewedBy,
		ProcessedBy:  row.ProcessedBy,
		Lines:        cloneReturnLines(row.Lines),
		RefundID:     row.RefundID,
		RefundStatus: row.RefundStatus,
		RefundAmount: row.RefundAmount,
		RequestedAt:  row.RequestedAt,
	}, nil
}

// Save writes the current ReturnRequest Active Record to its table.
func (request *ReturnRequest) Save() error {
	if request == nil || request.db == nil {
		return ErrDatabaseRequired
	}
	if request.ID == "" {
		return ErrReturnIDRequired
	}
	request.db.returns[request.ID] = returnRow{
		ID:           request.ID,
		OrderID:      request.OrderID,
		Status:       request.Status,
		Reason:       request.Reason,
		ReviewNote:   request.ReviewNote,
		RequestedBy:  request.RequestedBy,
		ReviewedBy:   request.ReviewedBy,
		ProcessedBy:  request.ProcessedBy,
		Lines:        cloneReturnLines(request.Lines),
		RefundID:     request.RefundID,
		RefundStatus: request.RefundStatus,
		RefundAmount: request.RefundAmount,
		RequestedAt:  request.RequestedAt,
	}
	return nil
}

func cloneReturnLines(lines []ReturnLine) []ReturnLine {
	if lines == nil {
		return nil
	}
	clone := make([]ReturnLine, len(lines))
	copy(clone, lines)
	return clone
}
