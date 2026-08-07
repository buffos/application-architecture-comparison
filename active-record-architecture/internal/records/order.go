package records

import (
	"errors"
	"strings"
	"time"
)

const (
	OrderStatusPendingReservation  = "PendingReservation"
	OrderStatusBackordered         = "Backordered"
	OrderStatusReadyForPayment     = "ReadyForPayment"
	OrderStatusPaymentReview       = "PaymentReview"
	OrderStatusReadyForFulfillment = "ReadyForFulfillment"
	OrderStatusPartiallyShipped    = "PartiallyShipped"
	OrderStatusShipped             = "Shipped"
	OrderStatusCancelled           = "Cancelled"
	PaymentStatusNotRequired       = "NotRequired"
	PaymentStatusPending           = "Pending"
	PaymentStatusAccepted          = "Accepted"
	PaymentStatusFailed            = "Failed"
	PaymentStatusManualReview      = "ManualReview"
	PaymentOutcomeAccept           = "accept"
	PaymentOutcomeFail             = "fail"
	PaymentOutcomeReview           = "review"
)

var (
	ErrOrderIDRequired            = errors.New("order id is required")
	ErrOrderNotFound              = errors.New("order not found")
	ErrOrderNotReservable         = errors.New("order is not awaiting reservation")
	ErrInsufficientStock          = errors.New("insufficient stock")
	ErrOrderNotPayable            = errors.New("order is not ready for payment")
	ErrPaymentOutcomeInvalid      = errors.New("payment outcome is invalid")
	ErrOrderNotShippable          = errors.New("order is not ready for fulfillment")
	ErrNoShipmentLines            = errors.New("order has no remaining shippable lines")
	ErrShippedByRequired          = errors.New("shipper is required")
	ErrOrderNotCancellable        = errors.New("order cannot be cancelled")
	ErrCancelledByRequired        = errors.New("cancelling actor is required")
	ErrCancellationReasonRequired = errors.New("cancellation reason is required")
	ErrStockReleaseInvalid        = errors.New("reserved stock cannot be released")
)

// OrderLine is a committed product snapshot embedded in an Order Active
// Record.
type OrderLine struct {
	ID                  string
	SKU                 string
	ProductNameSnapshot string
	ProductCategory     string
	OrderedQuantity     int
	ReservedQuantity    int
	ShippedQuantity     int
	ReturnedQuantity    int
	UnitPrice           int
	DiscountAmount      int
	ReturnWindowDays    int
	LineTotal           int
}

// Order is an Active Record for a committed commercial transaction.
type Order struct {
	db *Database

	ID                 string
	SourceQuoteID      string
	CustomerID         string
	Status             string
	RequestedBy        string
	PaymentID          string
	Lines              []OrderLine
	Total              int
	PaymentStatus      string
	ShippedAt          time.Time
	CancelledBy        string
	CancellationReason string
}

// RequestReturn validates shipped quantities and creates the passive return
// and refund records. It does not restock inventory or complete the refund.
func (order *Order) RequestReturn(lines []ReturnLine, reason string) (*ReturnRequest, error) {
	if order == nil || order.db == nil {
		return nil, ErrDatabaseRequired
	}
	if order.Status != OrderStatusShipped && order.Status != OrderStatusPartiallyShipped {
		return nil, ErrOrderNotReturnable
	}

	requestLines := lines
	if len(requestLines) == 0 {
		requestLines = make([]ReturnLine, 0, len(order.Lines))
		for _, orderLine := range order.Lines {
			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if remaining <= 0 {
				continue
			}
			requestLines = append(requestLines, ReturnLine{
				OrderLineID:     orderLine.ID,
				SKU:             orderLine.SKU,
				ProductCategory: orderLine.ProductCategory,
				Quantity:        remaining,
				UnitPrice:       orderLine.UnitPrice,
			})
		}
	}
	if len(requestLines) == 0 {
		return nil, ErrReturnLinesInvalid
	}

	normalizedLines := make([]ReturnLine, 0, len(requestLines))
	seenLines := make(map[string]struct{}, len(requestLines))
	refundAmount := 0
	for _, requestedLine := range requestLines {
		if requestedLine.Quantity <= 0 {
			return nil, ErrReturnLinesInvalid
		}
		if _, seen := seenLines[requestedLine.OrderLineID]; seen {
			return nil, ErrReturnLinesInvalid
		}
		seenLines[requestedLine.OrderLineID] = struct{}{}

		matched := false
		for _, orderLine := range order.Lines {
			if orderLine.ID != requestedLine.OrderLineID {
				continue
			}
			remaining := orderLine.ShippedQuantity - orderLine.ReturnedQuantity
			if requestedLine.Quantity > remaining {
				return nil, ErrReturnLinesInvalid
			}
			normalizedLines = append(normalizedLines, ReturnLine{
				OrderLineID:     orderLine.ID,
				SKU:             orderLine.SKU,
				ProductCategory: orderLine.ProductCategory,
				Quantity:        requestedLine.Quantity,
				UnitPrice:       orderLine.UnitPrice,
			})
			refundAmount += requestedLine.Quantity * orderLine.UnitPrice
			matched = true
			break
		}
		if !matched {
			return nil, ErrReturnLinesInvalid
		}
	}

	request := &ReturnRequest{
		db:           order.db,
		ID:           order.db.nextReturnID(),
		OrderID:      order.ID,
		Status:       ReturnStatusRequested,
		Reason:       reason,
		Lines:        cloneReturnLines(normalizedLines),
		RefundStatus: RefundStatusNotStarted,
		RefundAmount: refundAmount,
		RequestedAt:  time.Now(),
	}
	refund := &Refund{
		db:              order.db,
		ID:              order.db.nextRefundID(),
		ReturnRequestID: request.ID,
		OrderID:         order.ID,
		Amount:          refundAmount,
		Status:          RefundStatusNotStarted,
	}
	request.RefundID = refund.ID

	if err := request.Save(); err != nil {
		return nil, err
	}
	if err := refund.Save(); err != nil {
		return nil, err
	}
	return request, nil
}

// Cancel stops an unshipped order and releases every outstanding stock
// reservation. All stock rows are validated before any release is applied.
func (order *Order) Cancel(cancelledBy string, reason string) error {
	if order == nil || order.db == nil {
		return ErrDatabaseRequired
	}
	if cancelledBy == "" {
		return ErrCancelledByRequired
	}
	if reason == "" {
		return ErrCancellationReasonRequired
	}
	if order.Status != OrderStatusPendingReservation &&
		order.Status != OrderStatusBackordered &&
		order.Status != OrderStatusReadyForPayment &&
		order.Status != OrderStatusPaymentReview &&
		order.Status != OrderStatusReadyForFulfillment {
		return ErrOrderNotCancellable
	}

	for _, line := range order.Lines {
		if line.ReservedQuantity == 0 {
			continue
		}
		stock, err := FindStock(order.db, line.SKU)
		if err != nil || stock.Reserved < line.ReservedQuantity {
			return ErrStockReleaseInvalid
		}
	}

	for index := range order.Lines {
		line := &order.Lines[index]
		if line.ReservedQuantity == 0 {
			continue
		}
		stock, err := FindStock(order.db, line.SKU)
		if err != nil {
			return ErrStockReleaseInvalid
		}
		if err := stock.Release(line.ReservedQuantity); err != nil {
			return err
		}
		if err := stock.Save(); err != nil {
			return err
		}
		line.ReservedQuantity = 0
	}

	order.Status = OrderStatusCancelled
	order.CancelledBy = cancelledBy
	order.CancellationReason = reason
	return nil
}

// CapturePayment creates and persists a payment attempt, then updates the
// order lifecycle according to the simulated outcome.
func (order *Order) CapturePayment(outcome string) (*Payment, error) {
	if order == nil || order.db == nil {
		return nil, ErrDatabaseRequired
	}
	if order.Status != OrderStatusReadyForPayment {
		return nil, ErrOrderNotPayable
	}

	outcome = strings.ToLower(strings.TrimSpace(outcome))
	if outcome == "" {
		outcome = PaymentOutcomeAccept
	}
	paymentStatus := PaymentStatusAccepted
	orderStatus := OrderStatusReadyForFulfillment
	switch outcome {
	case PaymentOutcomeAccept:
		paymentStatus = PaymentStatusAccepted
	case PaymentOutcomeFail:
		paymentStatus = PaymentStatusFailed
		orderStatus = OrderStatusReadyForPayment
	case PaymentOutcomeReview:
		paymentStatus = PaymentStatusManualReview
		orderStatus = OrderStatusPaymentReview
	default:
		return nil, ErrPaymentOutcomeInvalid
	}

	payment := &Payment{
		db:      order.db,
		ID:      order.db.nextPaymentID(),
		OrderID: order.ID,
		Amount:  order.Total,
		Status:  paymentStatus,
	}
	if err := payment.Save(); err != nil {
		return nil, err
	}
	order.PaymentID = payment.ID
	order.PaymentStatus = paymentStatus
	order.Status = orderStatus
	return payment, nil
}

// CreateShipment creates and saves a full shipment, consumes the reserved
// stock rows, and updates the order's shipped quantities.
func (order *Order) CreateShipment(shippedBy string) (*Shipment, error) {
	if order == nil || order.db == nil {
		return nil, ErrDatabaseRequired
	}
	if shippedBy == "" {
		return nil, ErrShippedByRequired
	}
	if order.Status != OrderStatusReadyForFulfillment {
		return nil, ErrOrderNotShippable
	}

	shipmentLines := make([]ShipmentLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		remaining := line.ReservedQuantity - line.ShippedQuantity
		if remaining <= 0 {
			continue
		}
		shipmentLines = append(shipmentLines, ShipmentLine{OrderLineID: line.ID, SKU: line.SKU, Quantity: remaining})
	}
	if len(shipmentLines) == 0 {
		return nil, ErrNoShipmentLines
	}

	for _, shipmentLine := range shipmentLines {
		stock, err := FindStock(order.db, shipmentLine.SKU)
		if err != nil || stock.Reserved < shipmentLine.Quantity || stock.OnHand < shipmentLine.Quantity {
			return nil, ErrInsufficientStock
		}
	}

	shipment := &Shipment{
		db:        order.db,
		ID:        order.db.nextShipmentID(),
		OrderID:   order.ID,
		Status:    ShipmentStatusShipped,
		ShippedBy: shippedBy,
		ShippedAt: time.Now(),
		Lines:     cloneShipmentLines(shipmentLines),
	}

	for _, shipmentLine := range shipmentLines {
		for index := range order.Lines {
			if order.Lines[index].ID != shipmentLine.OrderLineID {
				continue
			}
			order.Lines[index].ShippedQuantity += shipmentLine.Quantity
			stock, err := FindStock(order.db, shipmentLine.SKU)
			if err != nil {
				return nil, err
			}
			if err := stock.Consume(shipmentLine.Quantity); err != nil {
				return nil, err
			}
			if err := stock.Save(); err != nil {
				return nil, err
			}
			break
		}
	}

	order.Status = OrderStatusShipped
	order.ShippedAt = shipment.ShippedAt
	if err := shipment.Save(); err != nil {
		return nil, err
	}
	return shipment, nil
}

// ReserveStock preflights every order line, then asks StockRecord Active
// Records to persist reservations. A shortage can be rejected or backordered
// according to the related Product Active Record.
func (order *Order) ReserveStock() error {
	if order == nil || order.db == nil {
		return ErrDatabaseRequired
	}
	if order.Status != OrderStatusPendingReservation {
		return ErrOrderNotReservable
	}

	type reservation struct {
		lineIndex int
		stock     *StockRecord
		quantity  int
	}
	reservations := make([]reservation, 0, len(order.Lines))
	planned := make(map[string]int)
	backordered := false

	for index, line := range order.Lines {
		if line.OrderedQuantity <= 0 {
			return ErrQuantityInvalid
		}
		stock, err := FindStock(order.db, line.SKU)
		available := 0
		if err == nil {
			available = stock.Available() - planned[line.SKU]
		}
		if err != nil || available < line.OrderedQuantity {
			policy := StockShortageRejectOrder
			if product, productErr := FindProduct(order.db, line.SKU); productErr == nil && product.StockShortagePolicy != "" {
				policy = product.StockShortagePolicy
			}
			if policy == StockShortageAllowBackorder {
				backordered = true
				continue
			}
			return ErrInsufficientStock
		}
		planned[line.SKU] += line.OrderedQuantity
		reservations = append(reservations, reservation{lineIndex: index, stock: stock, quantity: line.OrderedQuantity})
	}

	for _, item := range reservations {
		if err := item.stock.Reserve(item.quantity); err != nil {
			return err
		}
		if err := item.stock.Save(); err != nil {
			return err
		}
		order.Lines[item.lineIndex].ReservedQuantity = item.quantity
	}

	if backordered {
		order.Status = OrderStatusBackordered
		order.PaymentStatus = PaymentStatusNotRequired
	} else {
		order.Status = OrderStatusReadyForPayment
		order.PaymentStatus = PaymentStatusPending
	}
	return nil
}

// FindOrder loads an Order Active Record from the order table.
func FindOrder(db *Database, id string) (*Order, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrOrderIDRequired
	}

	row, ok := db.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return &Order{
		db:                 db,
		ID:                 row.ID,
		SourceQuoteID:      row.SourceQuoteID,
		CustomerID:         row.CustomerID,
		Status:             row.Status,
		RequestedBy:        row.RequestedBy,
		PaymentID:          row.PaymentID,
		PaymentStatus:      row.PaymentStatus,
		ShippedAt:          row.ShippedAt,
		CancelledBy:        row.CancelledBy,
		CancellationReason: row.CancellationReason,
		Lines:              cloneOrderLines(row.Lines),
		Total:              row.Total,
	}, nil
}

// Save writes the current Order Active Record to its table.
func (order *Order) Save() error {
	if order == nil || order.db == nil {
		return ErrDatabaseRequired
	}
	if order.ID == "" {
		return ErrOrderIDRequired
	}

	order.db.orders[order.ID] = orderRow{
		ID:                 order.ID,
		SourceQuoteID:      order.SourceQuoteID,
		CustomerID:         order.CustomerID,
		Status:             order.Status,
		RequestedBy:        order.RequestedBy,
		PaymentID:          order.PaymentID,
		PaymentStatus:      order.PaymentStatus,
		ShippedAt:          order.ShippedAt,
		CancelledBy:        order.CancelledBy,
		CancellationReason: order.CancellationReason,
		Lines:              cloneOrderLines(order.Lines),
		Total:              order.Total,
	}
	return nil
}

func cloneOrderLines(lines []OrderLine) []OrderLine {
	if lines == nil {
		return nil
	}
	clone := make([]OrderLine, len(lines))
	copy(clone, lines)
	return clone
}
