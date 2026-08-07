package records

import (
	"errors"
	"time"
)

const (
	ReturnStatusRequested  = "Requested"
	ReturnStatusRefunded   = "Refunded"
	RefundStatusNotStarted = "NotStarted"
	RefundStatusCompleted  = "Completed"
)

var (
	ErrOrderNotReturnable  = errors.New("order is not returnable")
	ErrReturnIDRequired    = errors.New("return id is required")
	ErrReturnNotFound      = errors.New("return request not found")
	ErrReturnLinesInvalid  = errors.New("return lines are invalid")
	ErrReturnNotAcceptable = errors.New("return is not requested")
	ErrReturnNotRefundable = errors.New("return refund cannot be completed")
	ErrReturnOrderMissing  = errors.New("return order not found")
	ErrReturnStockMissing  = errors.New("return stock record not found")
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

// Accept applies the first, intentionally compressed return workflow: it
// marks order quantities returned, restocks inventory, and completes the
// linked refund.
func (request *ReturnRequest) Accept() error {
	if request == nil || request.db == nil {
		return ErrDatabaseRequired
	}
	if request.Status != ReturnStatusRequested {
		return ErrReturnNotAcceptable
	}

	order, err := FindOrder(request.db, request.OrderID)
	if err != nil {
		return ErrReturnOrderMissing
	}
	refund, err := FindRefund(request.db, request.RefundID)
	if err != nil {
		return ErrReturnNotRefundable
	}
	if refund.Status != RefundStatusNotStarted {
		return ErrReturnNotRefundable
	}

	type validatedLine struct {
		orderLineIndex int
		stock          *StockRecord
		quantity       int
	}
	validated := make([]validatedLine, 0, len(request.Lines))
	for _, returnLine := range request.Lines {
		if returnLine.Quantity <= 0 {
			return ErrReturnLinesInvalid
		}

		orderLineIndex := -1
		for index, orderLine := range order.Lines {
			if orderLine.ID != returnLine.OrderLineID {
				continue
			}
			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if returnLine.Quantity > remaining {
				return ErrReturnLinesInvalid
			}
			orderLineIndex = index
			break
		}
		if orderLineIndex < 0 {
			return ErrReturnLinesInvalid
		}

		stock, err := FindStock(request.db, order.Lines[orderLineIndex].SKU)
		if err != nil {
			return ErrReturnStockMissing
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
			return err
		}
	}

	refund.Status = RefundStatusCompleted
	if err := refund.Save(); err != nil {
		return err
	}
	request.Status = ReturnStatusRefunded
	request.RefundStatus = RefundStatusCompleted
	if err := order.Save(); err != nil {
		return err
	}
	return request.Save()
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
