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
	ErrActorRequired           = errors.New("actor is required")
	ErrReturnSelectionInvalid  = errors.New("return selection is invalid")
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

type ReturnSelection struct {
	ProductSKU ProductSKU
	Quantity   int
}

func (line ReturnLine) ProductSKU() ProductSKU           { return line.sku }
func (line ReturnLine) ProductCategory() ProductCategory { return line.category }
func (line ReturnLine) Quantity() int                    { return line.quantity }
func (line ReturnLine) UnitPrice() Money                 { return line.unitPrice }

// ReturnRequest is the aggregate root for return intent and review state.
type ReturnRequest struct {
	id          ReturnRequestID
	orderID     OrderID
	customerID  CustomerID
	reason      string
	status      ReturnStatus
	lines       []ReturnLine
	requestedBy string
	reviewedBy  string
	processedBy string
}

func NewReturnRequestFromShippedOrder(id ReturnRequestID, order ordering.Order, reason string) (ReturnRequest, error) {
	return NewReturnRequestFromOrderSelection(id, order, reason, nil)
}

func NewReturnRequestFromOrderSelection(id ReturnRequestID, order ordering.Order, reason string, selections []ReturnSelection) (ReturnRequest, error) {
	if id == "" {
		return ReturnRequest{}, ErrReturnRequestIDRequired
	}
	if order.Status() != ordering.OrderStatusShipped {
		return ReturnRequest{}, ErrOrderNotShipped
	}
	if reason == "" {
		return ReturnRequest{}, ErrReturnReasonRequired
	}
	if len(selections) == 0 {
		for _, line := range order.Lines() {
			quantity := line.ShippedQuantity()
			if quantity == 0 {
				quantity = line.Quantity()
			}
			selections = append(selections, ReturnSelection{ProductSKU: ProductSKU(line.ProductSKU()), Quantity: quantity})
		}
	}
	lines := make([]ReturnLine, 0, len(selections))
	for _, selection := range selections {
		if selection.Quantity <= 0 {
			return ReturnRequest{}, ErrReturnSelectionInvalid
		}
		matched := false
		for _, line := range order.Lines() {
			if ProductSKU(line.ProductSKU()) != selection.ProductSKU {
				continue
			}
			shippedQuantity := line.ShippedQuantity()
			if shippedQuantity == 0 {
				shippedQuantity = line.Quantity()
			}
			if selection.Quantity > shippedQuantity {
				return ReturnRequest{}, ErrReturnSelectionInvalid
			}
			price, err := NewMoney(line.UnitPrice().Cents(), line.UnitPrice().Currency())
			if err != nil {
				return ReturnRequest{}, err
			}
			lines = append(lines, ReturnLine{sku: ProductSKU(line.ProductSKU()), category: ProductCategory(line.ProductCategory()), quantity: selection.Quantity, unitPrice: price})
			matched = true
			break
		}
		if !matched {
			return ReturnRequest{}, ErrReturnSelectionInvalid
		}
	}
	return ReturnRequest{id: id, orderID: OrderID(order.ID()), customerID: CustomerID(order.CustomerID()), reason: reason, status: ReturnStatusRequested, lines: lines}, nil
}

func (r ReturnRequest) ID() ReturnRequestID    { return r.id }
func (r ReturnRequest) OrderID() OrderID       { return r.orderID }
func (r ReturnRequest) CustomerID() CustomerID { return r.customerID }
func (r ReturnRequest) Reason() string         { return r.reason }
func (r ReturnRequest) Status() ReturnStatus   { return r.status }
func (r ReturnRequest) Lines() []ReturnLine    { return append([]ReturnLine(nil), r.lines...) }
func (r ReturnRequest) RequestedBy() string    { return r.requestedBy }
func (r ReturnRequest) ReviewedBy() string     { return r.reviewedBy }
func (r ReturnRequest) ProcessedBy() string    { return r.processedBy }

func (r *ReturnRequest) AssignRequester(actor string) error {
	if actor == "" {
		return ErrActorRequired
	}
	r.requestedBy = actor
	return nil
}

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

func (r *ReturnRequest) ReviewBy(decision ReviewDecision, actor string) error {
	if actor == "" {
		return ErrActorRequired
	}
	if err := r.Review(decision); err != nil {
		return err
	}
	r.reviewedBy = actor
	return nil
}

func (r *ReturnRequest) ProcessBy(actor string) error {
	if actor == "" {
		return ErrActorRequired
	}
	if r.status != ReturnStatusAccepted {
		return ErrReturnNotReviewable
	}
	r.processedBy = actor
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
