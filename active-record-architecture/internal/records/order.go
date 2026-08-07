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
	PaymentReviewDecisionAccept    = "accept"
	PaymentReviewDecisionReject    = "reject"
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
	ErrShipmentLinesInvalid       = errors.New("shipment lines are invalid")
	ErrShippedByRequired          = errors.New("shipper is required")
	ErrOrderNotCancellable        = errors.New("order cannot be cancelled")
	ErrCancelledByRequired        = errors.New("cancelling actor is required")
	ErrCancellationReasonRequired = errors.New("cancellation reason is required")
	ErrStockReleaseInvalid        = errors.New("reserved stock cannot be released")
	ErrOrderNotInPaymentReview    = errors.New("order is not in payment review")
	ErrPaymentReviewMissing       = errors.New("payment review record not found")
	ErrPaymentReviewerRequired    = errors.New("payment reviewer is required")
	ErrPaymentDecisionInvalid     = errors.New("payment review decision is invalid")
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
func (order *Order) RequestReturn(lines []ReturnLine, reason string, requestedBy string) (*ReturnRequest, error) {
	return order.RequestReturnAt(lines, reason, requestedBy, time.Now())
}

// RequestReturnAt is the deterministic form of RequestReturn used by tests
// and demonstrations.
func (order *Order) RequestReturnAt(lines []ReturnLine, reason string, requestedBy string, requestedAt time.Time) (*ReturnRequest, error) {
	if order == nil || order.db == nil {
		return nil, ErrDatabaseRequired
	}
	if requestedBy == "" {
		return nil, ErrActorRequired
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
				OrderLineID:      orderLine.ID,
				SKU:              orderLine.SKU,
				ProductCategory:  orderLine.ProductCategory,
				Quantity:         remaining,
				UnitPrice:        orderLine.UnitPrice,
				ReturnWindowDays: orderLine.ReturnWindowDays,
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
				OrderLineID:      orderLine.ID,
				SKU:              orderLine.SKU,
				ProductCategory:  orderLine.ProductCategory,
				Quantity:         requestedLine.Quantity,
				UnitPrice:        orderLine.UnitPrice,
				ReturnWindowDays: orderLine.ReturnWindowDays,
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
		RequestedBy:  requestedBy,
		Lines:        cloneReturnLines(normalizedLines),
		RefundStatus: RefundStatusNotStarted,
		RefundAmount: refundAmount,
		RequestedAt:  requestedAt,
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

// ResolvePaymentReview records a manual payment decision and advances the
// order to fulfillment or payment retry.
func (order *Order) ResolvePaymentReview(reviewedBy string, decision string, comment string) error {
	if order == nil || order.db == nil {
		return ErrDatabaseRequired
	}
	if reviewedBy == "" {
		return ErrPaymentReviewerRequired
	}
	if order.Status != OrderStatusPaymentReview {
		return ErrOrderNotInPaymentReview
	}

	payment, err := FindPayment(order.db, order.PaymentID)
	if err != nil || payment.Status != PaymentStatusManualReview {
		return ErrPaymentReviewMissing
	}

	decision = strings.ToLower(strings.TrimSpace(decision))
	nextPaymentStatus := ""
	nextOrderStatus := ""
	switch decision {
	case PaymentReviewDecisionAccept:
		nextPaymentStatus = PaymentStatusAccepted
		nextOrderStatus = OrderStatusReadyForFulfillment
	case PaymentReviewDecisionReject:
		nextPaymentStatus = PaymentStatusFailed
		nextOrderStatus = OrderStatusReadyForPayment
	default:
		return ErrPaymentDecisionInvalid
	}

	payment.Status = nextPaymentStatus
	payment.ReviewedBy = reviewedBy
	payment.DecisionComment = comment
	if err := payment.Save(); err != nil {
		return err
	}
	order.Status = nextOrderStatus
	order.PaymentStatus = nextPaymentStatus
	return nil
}

// CreateShipment creates a full shipment by delegating to the partial-shipment
// operation with no explicit line selection.
func (order *Order) CreateShipment(shippedBy string) (*Shipment, error) {
	return order.CreatePartialShipment(shippedBy, nil)
}

// CreatePartialShipment creates a shipment for selected remaining quantities,
// consumes the matching reservations, and derives the aggregate order status.
// An empty selection means all remaining quantities.
func (order *Order) CreatePartialShipment(shippedBy string, requestedLines []ShipmentLine) (*Shipment, error) {
	if order == nil || order.db == nil {
		return nil, ErrDatabaseRequired
	}
	if shippedBy == "" {
		return nil, ErrShippedByRequired
	}
	if order.Status != OrderStatusReadyForFulfillment && order.Status != OrderStatusPartiallyShipped {
		return nil, ErrOrderNotShippable
	}

	if len(requestedLines) == 0 {
		requestedLines = make([]ShipmentLine, 0, len(order.Lines))
		for _, line := range order.Lines {
			remaining := line.ReservedQuantity - line.ShippedQuantity
			if remaining <= 0 {
				continue
			}
			requestedLines = append(requestedLines, ShipmentLine{OrderLineID: line.ID, SKU: line.SKU, Quantity: remaining})
		}
	}
	if len(requestedLines) == 0 {
		return nil, ErrNoShipmentLines
	}

	type selectedLine struct {
		orderLineIndex int
		shipmentLine   ShipmentLine
	}
	selected := make([]selectedLine, 0, len(requestedLines))
	selectedByOrderLine := make(map[int]int, len(requestedLines))
	plannedBySKU := make(map[string]int)
	for _, requestedLine := range requestedLines {
		if requestedLine.OrderLineID == "" || requestedLine.Quantity <= 0 {
			return nil, ErrShipmentLinesInvalid
		}
		orderLineIndex := -1
		for index, orderLine := range order.Lines {
			if orderLine.ID == requestedLine.OrderLineID {
				orderLineIndex = index
				break
			}
		}
		if orderLineIndex < 0 || selectedByOrderLine[orderLineIndex] > 0 {
			return nil, ErrShipmentLinesInvalid
		}
		remaining := order.Lines[orderLineIndex].ReservedQuantity - order.Lines[orderLineIndex].ShippedQuantity
		if requestedLine.Quantity > remaining {
			return nil, ErrShipmentLinesInvalid
		}
		canonicalLine := ShipmentLine{
			OrderLineID: requestedLine.OrderLineID,
			SKU:         order.Lines[orderLineIndex].SKU,
			Quantity:    requestedLine.Quantity,
		}
		selected = append(selected, selectedLine{orderLineIndex: orderLineIndex, shipmentLine: canonicalLine})
		selectedByOrderLine[orderLineIndex] = requestedLine.Quantity
		plannedBySKU[canonicalLine.SKU] += requestedLine.Quantity
	}

	stockBySKU := make(map[string]*StockRecord, len(plannedBySKU))
	for sku, quantity := range plannedBySKU {
		stock, err := FindStock(order.db, sku)
		if err != nil || stock.Reserved < quantity || stock.OnHand < quantity {
			return nil, ErrInsufficientStock
		}
		stockBySKU[sku] = stock
	}

	shipmentLines := make([]ShipmentLine, 0, len(selected))
	for _, item := range selected {
		shipmentLines = append(shipmentLines, item.shipmentLine)
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

	for _, item := range selected {
		order.Lines[item.orderLineIndex].ShippedQuantity += item.shipmentLine.Quantity
	}
	for sku, quantity := range plannedBySKU {
		stock := stockBySKU[sku]
		if err := stock.Consume(quantity); err != nil {
			return nil, err
		}
		if err := stock.Save(); err != nil {
			return nil, err
		}
	}

	fullyShipped := true
	for _, line := range order.Lines {
		remaining := line.ReservedQuantity - line.ShippedQuantity
		if remaining > 0 {
			fullyShipped = false
			break
		}
	}
	if fullyShipped {
		order.Status = OrderStatusShipped
	} else {
		order.Status = OrderStatusPartiallyShipped
	}
	if order.ShippedAt.IsZero() {
		order.ShippedAt = shipment.ShippedAt
	}
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
