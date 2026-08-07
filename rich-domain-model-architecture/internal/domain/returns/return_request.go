package returns

import (
	"errors"

	"rich-domain-model-architecture/internal/domain/ordering"
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
type ProductCategory string

const (
	ProductCategoryStandard    ProductCategory = "Standard"
	ProductCategoryCustomBuild ProductCategory = "CustomBuild"
	ProductCategoryClearance   ProductCategory = "Clearance"
)

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

// ReturnRequest owns return intent and review state.
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

	selected := make(map[ProductSKU]int)
	if len(selections) == 0 {
		for _, line := range order.Lines() {
			quantity := line.ShippedQuantity()
			if quantity == 0 {
				quantity = line.Quantity()
			}
			selected[ProductSKU(line.ProductSKU())] = quantity
		}
	} else {
		for _, selection := range selections {
			if selection.Quantity <= 0 {
				return ReturnRequest{}, ErrReturnSelectionInvalid
			}
			selected[selection.ProductSKU] += selection.Quantity
		}
	}

	lines := make([]ReturnLine, 0, len(selected))
	matched := make(map[ProductSKU]bool, len(selected))
	for _, line := range order.Lines() {
		sku := ProductSKU(line.ProductSKU())
		quantity, ok := selected[sku]
		if !ok {
			continue
		}
		shippedQuantity := line.ShippedQuantity()
		if shippedQuantity == 0 {
			shippedQuantity = line.Quantity()
		}
		if quantity <= 0 || quantity > shippedQuantity {
			return ReturnRequest{}, ErrReturnSelectionInvalid
		}
		price, err := NewMoney(line.UnitPrice().Cents(), line.UnitPrice().Currency())
		if err != nil {
			return ReturnRequest{}, err
		}
		lines = append(lines, ReturnLine{
			sku:       sku,
			category:  ProductCategory(line.ProductCategory()),
			quantity:  quantity,
			unitPrice: price,
		})
		matched[sku] = true
	}
	for sku := range selected {
		if !matched[sku] {
			return ReturnRequest{}, ErrReturnSelectionInvalid
		}
	}
	return ReturnRequest{
		id:         id,
		orderID:    OrderID(order.ID()),
		customerID: CustomerID(order.CustomerID()),
		reason:     reason,
		status:     ReturnStatusRequested,
		lines:      lines,
	}, nil
}

func (request ReturnRequest) ID() ReturnRequestID    { return request.id }
func (request ReturnRequest) OrderID() OrderID       { return request.orderID }
func (request ReturnRequest) CustomerID() CustomerID { return request.customerID }
func (request ReturnRequest) Reason() string         { return request.reason }
func (request ReturnRequest) Status() ReturnStatus   { return request.status }
func (request ReturnRequest) Lines() []ReturnLine    { return append([]ReturnLine(nil), request.lines...) }
func (request ReturnRequest) RequestedBy() string    { return request.requestedBy }
func (request ReturnRequest) ReviewedBy() string     { return request.reviewedBy }
func (request ReturnRequest) ProcessedBy() string    { return request.processedBy }

func (request *ReturnRequest) Review(decision ReviewDecision) error {
	if request.status != ReturnStatusRequested {
		return ErrReturnNotReviewable
	}
	switch decision {
	case ReviewDecisionAccept:
		request.status = ReturnStatusAccepted
	case ReviewDecisionReject:
		request.status = ReturnStatusRejected
	default:
		return ErrReviewDecisionInvalid
	}
	return nil
}

func (request *ReturnRequest) Accept() error {
	return request.Review(ReviewDecisionAccept)
}

func (request *ReturnRequest) Reject() error {
	return request.Review(ReviewDecisionReject)
}

func (request *ReturnRequest) AssignRequester(actor string) error {
	if actor == "" {
		return ErrActorRequired
	}
	request.requestedBy = actor
	return nil
}

func (request *ReturnRequest) ReviewBy(decision ReviewDecision, actor string) error {
	if actor == "" {
		return ErrActorRequired
	}
	if err := request.Review(decision); err != nil {
		return err
	}
	request.reviewedBy = actor
	return nil
}

func (request *ReturnRequest) ProcessBy(actor string) error {
	if actor == "" {
		return ErrActorRequired
	}
	if request.status != ReturnStatusAccepted {
		return ErrReturnNotReviewable
	}
	request.processedBy = actor
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

func (refund Refund) ID() RefundID                     { return refund.id }
func (refund Refund) ReturnRequestID() ReturnRequestID { return refund.returnRequestID }
func (refund Refund) Amount() Money                    { return refund.amount }
func (refund Refund) Status() RefundStatus             { return refund.status }

func (refund *Refund) Issue() error {
	if refund.status != RefundStatusPending {
		return ErrRefundNotIssuable
	}
	refund.status = RefundStatusIssued
	return nil
}
