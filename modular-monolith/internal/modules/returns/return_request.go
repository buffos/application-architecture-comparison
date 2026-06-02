package returns

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrReturnRequestNotFound = errors.New("return request not found")

const ReturnRequestStatusRefunded = "Refunded"

var returnRequestSequence uint64

type ReturnRequest struct {
	ID         string
	OrderID    string
	CustomerID string
	Reason     string
	Status     string
	Lines      []ReturnRequestLine
}

type ReturnRequestLine struct {
	ProductSKU      string
	ProductName     string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

func NewRefundedReturnRequest(order ReturnableOrder, reason string) ReturnRequest {
	id := atomic.AddUint64(&returnRequestSequence, 1)

	lines := make([]ReturnRequestLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ReturnRequestLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
		})
	}

	return ReturnRequest{
		ID:         fmt.Sprintf("return-%03d", id),
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Reason:     reason,
		Status:     ReturnRequestStatusRefunded,
		Lines:      lines,
	}
}
