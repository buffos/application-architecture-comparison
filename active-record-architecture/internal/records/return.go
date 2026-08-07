package records

import (
	"errors"
	"time"
)

const (
	ReturnStatusRequested  = "Requested"
	RefundStatusNotStarted = "NotStarted"
)

var (
	ErrOrderNotReturnable = errors.New("order is not returnable")
	ErrReturnIDRequired   = errors.New("return id is required")
	ErrReturnNotFound     = errors.New("return request not found")
	ErrReturnLinesInvalid = errors.New("return lines are invalid")
)

// ReturnLine is a passive request for a quantity from an order-line snapshot.
type ReturnLine struct {
	OrderLineID     string
	SKU             string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

// ReturnRequest is an Active Record for the requested reverse flow. Its
// business transitions are added in later lessons.
type ReturnRequest struct {
	db *Database

	ID           string
	OrderID      string
	Status       string
	Reason       string
	Lines        []ReturnLine
	RefundID     string
	RefundStatus string
	RefundAmount int
	RequestedAt  time.Time
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
