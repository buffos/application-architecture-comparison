package orders

import (
	"fmt"

	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/payments"
	"component-based-architecture/internal/components/quotes"
	"component-based-architecture/internal/components/shipments"
)

// Component owns order creation and order state for this lesson.
type Component struct {
	quotes    quotes.ApprovedQuoteSource
	stock     inventory.StockKeeper
	payments  payments.Processor
	shipments shipments.Creator
	orders    map[string]Order
	nextID    int
}

func NewComponent(quotes quotes.ApprovedQuoteSource, stock inventory.StockKeeper, payments payments.Processor, shipments shipments.Creator) *Component {
	return &Component{quotes: quotes, stock: stock, payments: payments, shipments: shipments, orders: make(map[string]Order)}
}

type ConvertQuoteToOrderCommand struct {
	QuoteID string
}

type ConvertQuoteToOrderResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CapturePaymentCommand struct {
	OrderID string
}

type CapturePaymentResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ApprovePaymentReviewCommand struct{ OrderID string }

type ApprovePaymentReviewResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CreateShipmentCommand struct {
	OrderID string
	Lines   []ShipmentSelection
}

type ShipmentSelection struct {
	ProductSKU string
	Quantity   int
}

type CreateShipmentResult struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CancelOrderCommand struct {
	OrderID string
}

type CancelOrderResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

var _ Reader = (*Component)(nil)

func (c *Component) ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error) {
	quote, err := c.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	reservations := make([]inventory.ReservationItem, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		reservations = append(reservations, inventory.ReservationItem{ProductSKU: line.ProductSKU, Quantity: line.Quantity})
	}
	if err := c.stock.Reserve(reservations); err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	c.nextID++
	order := newOrderFromApprovedQuote(fmt.Sprintf("order-%03d", c.nextID), quote)
	c.orders[order.ID] = order

	return ConvertQuoteToOrderResult{
		OrderID: order.ID, QuoteID: order.QuoteID, CustomerID: order.CustomerID,
		Status: order.Status, LineCount: len(order.Lines),
	}, nil
}

func (c *Component) CapturePayment(command CapturePaymentCommand) (CapturePaymentResult, error) {
	order, ok := c.orders[command.OrderID]
	if !ok {
		return CapturePaymentResult{}, ErrOrderNotFound
	}

	amount := 0
	for _, line := range order.Lines {
		amount += line.Quantity * line.UnitPrice
	}
	capture, err := c.payments.Capture(payments.PaymentRequest{OrderID: order.ID, CustomerID: order.CustomerID, Amount: amount})
	if err != nil {
		return CapturePaymentResult{}, err
	}
	if capture.Outcome == payments.CaptureOutcomeReview {
		if err := order.MarkPaymentReview(); err != nil {
			return CapturePaymentResult{}, err
		}
	} else if err := order.MarkPaid(); err != nil {
		return CapturePaymentResult{}, err
	}
	c.orders[order.ID] = order

	return CapturePaymentResult{OrderID: order.ID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(order.Lines)}, nil
}

func (c *Component) ApprovePaymentReview(command ApprovePaymentReviewCommand) (ApprovePaymentReviewResult, error) {
	order, ok := c.orders[command.OrderID]
	if !ok {
		return ApprovePaymentReviewResult{}, ErrOrderNotFound
	}
	if err := order.ApprovePaymentReview(); err != nil {
		return ApprovePaymentReviewResult{}, err
	}
	c.orders[order.ID] = order
	return ApprovePaymentReviewResult{OrderID: order.ID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(order.Lines)}, nil
}

func (c *Component) CreateShipment(command CreateShipmentCommand) (CreateShipmentResult, error) {
	order, ok := c.orders[command.OrderID]
	if !ok {
		return CreateShipmentResult{}, ErrOrderNotFound
	}
	if order.Status != OrderStatusPaid && order.Status != OrderStatusPartiallyShipped {
		return CreateShipmentResult{}, ErrOrderNotShippable
	}

	selections := append([]ShipmentSelection(nil), command.Lines...)
	if len(selections) == 0 {
		for _, line := range order.Lines {
			if remaining := line.Quantity - line.ShippedQuantity; remaining > 0 {
				selections = append(selections, ShipmentSelection{ProductSKU: line.ProductSKU, Quantity: remaining})
			}
		}
	}
	lines := make([]shipments.ShipmentLine, 0, len(order.Lines))
	for _, selection := range selections {
		for _, line := range order.Lines {
			if line.ProductSKU == selection.ProductSKU {
				lines = append(lines, shipments.ShipmentLine{ProductSKU: line.ProductSKU, ProductName: line.ProductName, Quantity: selection.Quantity})
				break
			}
		}
	}
	if len(lines) == 0 {
		return CreateShipmentResult{}, ErrOrderNotShippable
	}
	shipment, err := c.shipments.Create(shipments.ShipmentRequest{OrderID: order.ID, CustomerID: order.CustomerID, Lines: lines})
	if err != nil {
		return CreateShipmentResult{}, err
	}
	if err := order.ApplyShipment(selections); err != nil {
		return CreateShipmentResult{}, err
	}
	if order.Status == OrderStatusShipped {
		order.ShippedAt = shipment.ShippedAt
	}
	c.orders[order.ID] = order

	return CreateShipmentResult{ShipmentID: shipment.ID, OrderID: order.ID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(order.Lines)}, nil
}

func (c *Component) CancelOrder(command CancelOrderCommand) (CancelOrderResult, error) {
	order, ok := c.orders[command.OrderID]
	if !ok {
		return CancelOrderResult{}, ErrOrderNotFound
	}
	if err := order.Cancel(); err != nil {
		return CancelOrderResult{}, err
	}

	releases := make([]inventory.ReleaseItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		releases = append(releases, inventory.ReleaseItem{ProductSKU: line.ProductSKU, Quantity: line.Quantity})
	}
	if err := c.stock.Release(releases); err != nil {
		return CancelOrderResult{}, err
	}
	c.orders[order.ID] = order

	return CancelOrderResult{OrderID: order.ID, CustomerID: order.CustomerID, Status: order.Status, LineCount: len(order.Lines)}, nil
}
