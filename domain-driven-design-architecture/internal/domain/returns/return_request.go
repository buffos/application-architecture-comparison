package returns

import (
	"errors"

	"domain-driven-design-architecture/internal/domain/ordering"
)

var (
	ErrReturnRequestIDRequired = errors.New("return request id is required")
	ErrReturnReasonRequired    = errors.New("return reason is required")
	ErrOrderNotShipped         = errors.New("order is not shipped")
	ErrReturnNotReviewable     = errors.New("return request is not reviewable")
	ErrReviewDecisionInvalid   = errors.New("review decision is invalid")
	ErrRefundIDRequired        = errors.New("refund id is required")
	ErrRefundNotIssuable       = errors.New("refund is not issuable")
)

type ReturnRequestID string
type OrderID string
type CustomerID string
type RefundID string
type ProductSKU string

type ReturnStatus string

const (
	ReturnStatusRequested ReturnStatus = "Requested"
	ReturnStatusAccepted  ReturnStatus = "Accepted"
	ReturnStatusRejected  ReturnStatus = "Rejected"
)

type ReviewDecision string

const (
	ReviewDecisionAccept ReviewDecision = "Accept"
	ReviewDecisionReject ReviewDecision = "Reject"
)

type ReturnLine struct {
	sku       ProductSKU
	category  ProductCategory
	quantity  int
	unitPrice Money
}

func (line ReturnLine) ProductSKU() ProductSKU           { return line.sku }
func (line ReturnLine) ProductCategory() ProductCategory { return line.category }
func (line ReturnLine) Quantity() int                    { return line.quantity }
func (line ReturnLine) UnitPrice() Money                 { return line.unitPrice }

// ReturnRequest is the aggregate root for return intent and review state.
type ReturnRequest struct {
	id         ReturnRequestID
	orderID    OrderID
	customerID CustomerID
	reason     string
	status     ReturnStatus
	lines      []ReturnLine
}

func NewReturnRequestFromShippedOrder(id ReturnRequestID, order ordering.Order, reason string) (ReturnRequest, error) {
	if id == "" {
		return ReturnRequest{}, ErrReturnRequestIDRequired
	}
	if order.Status() != ordering.OrderStatusShipped {
		return ReturnRequest{}, ErrOrderNotShipped
	}
	if reason == "" {
		return ReturnRequest{}, ErrReturnReasonRequired
	}
	lines := make([]ReturnLine, 0, len(order.Lines()))
	for _, line := range order.Lines() {
		price, err := NewMoney(line.UnitPrice().Cents(), line.UnitPrice().Currency())
		if err != nil {
			return ReturnRequest{}, err
		}
		lines = append(lines, ReturnLine{sku: ProductSKU(line.ProductSKU()), category: ProductCategory(line.ProductCategory()), quantity: line.Quantity(), unitPrice: price})
	}
	return ReturnRequest{id: id, orderID: OrderID(order.ID()), customerID: CustomerID(order.CustomerID()), reason: reason, status: ReturnStatusRequested, lines: lines}, nil
}

func (r ReturnRequest) ID() ReturnRequestID    { return r.id }
func (r ReturnRequest) OrderID() OrderID       { return r.orderID }
func (r ReturnRequest) CustomerID() CustomerID { return r.customerID }
func (r ReturnRequest) Reason() string         { return r.reason }
func (r ReturnRequest) Status() ReturnStatus   { return r.status }
func (r ReturnRequest) Lines() []ReturnLine    { return append([]ReturnLine(nil), r.lines...) }

func (r *ReturnRequest) Accept() error {
	return r.Review(ReviewDecisionAccept)
}

func (r *ReturnRequest) Reject() error {
	return r.Review(ReviewDecisionReject)
}

func (r *ReturnRequest) Review(decision ReviewDecision) error {
	if r.status != ReturnStatusRequested {
		return ErrReturnNotReviewable
	}
	switch decision {
	case ReviewDecisionAccept:
		r.status = ReturnStatusAccepted
	case ReviewDecisionReject:
		r.status = ReturnStatusRejected
	default:
		return ErrReviewDecisionInvalid
	}
	return nil
}

type RefundStatus string

const (
	RefundStatusPending RefundStatus = "Pending"
	RefundStatusIssued  RefundStatus = "Issued"
)

type Refund struct {
	id              RefundID
	returnRequestID ReturnRequestID
	amount          Money
	status          RefundStatus
}

func NewRefund(id RefundID, request ReturnRequestID, amount Money) (Refund, error) {
	if id == "" {
		return Refund{}, ErrRefundIDRequired
	}
	if request == "" {
		return Refund{}, ErrReturnRequestIDRequired
	}
	if amount.Currency() == "" {
		return Refund{}, ErrCurrencyRequired
	}
	return Refund{id: id, returnRequestID: request, amount: amount, status: RefundStatusPending}, nil
}

func (r Refund) ID() RefundID                     { return r.id }
func (r Refund) ReturnRequestID() ReturnRequestID { return r.returnRequestID }
func (r Refund) Amount() Money                    { return r.amount }
func (r Refund) Status() RefundStatus             { return r.status }

func (r *Refund) Issue() error {
	if r.status != RefundStatusPending {
		return ErrRefundNotIssuable
	}
	r.status = RefundStatusIssued
	return nil
}
