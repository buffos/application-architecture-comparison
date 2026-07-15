package orders

import "microkernel-architecture/internal/kernel"

type Service struct {
	orders  Repository
	quotes  kernel.ApprovedQuoteProvider
	stock   kernel.InventoryReservation
	release kernel.InventoryRelease
	pay     kernel.PaymentCapture
	ship    kernel.ShipmentCreation
	clock   kernel.Clock
}

func NewService(orders Repository, quotes kernel.ApprovedQuoteProvider, stock kernel.InventoryReservation, release kernel.InventoryRelease, pay kernel.PaymentCapture, ship kernel.ShipmentCreation, clock kernel.Clock) Service {
	return Service{
		orders:  orders,
		quotes:  quotes,
		stock:   stock,
		release: release,
		pay:     pay,
		ship:    ship,
		clock:   clock,
	}
}

func (s Service) ConvertQuoteToOrder(command kernel.ConvertQuoteToOrderCommand) (kernel.ConvertQuoteToOrderResult, error) {
	quote, err := s.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

	lines := make([]OrderLine, 0, len(quote.Lines))
	reservationItems := make([]kernel.InventoryReservationItem, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:       line.ProductSKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
		reservationItems = append(reservationItems, kernel.InventoryReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	order := NewOrderFromApprovedQuote(quote.QuoteID, quote.CustomerID, lines)
	if err := s.stock.Reserve(reservationItems); err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

	return kernel.ConvertQuoteToOrderResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) CapturePayment(command kernel.CapturePaymentCommand) (kernel.CapturePaymentResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return kernel.CapturePaymentResult{}, err
	}

	if err := s.pay.Capture(order.ID, order.TotalAmount()); err != nil {
		return kernel.CapturePaymentResult{}, err
	}

	if err := order.MarkPaid(); err != nil {
		return kernel.CapturePaymentResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return kernel.CapturePaymentResult{}, err
	}

	return kernel.CapturePaymentResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) CreateShipment(command kernel.CreateShipmentCommand) (kernel.CreateShipmentResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return kernel.CreateShipmentResult{}, err
	}

	lines := make([]kernel.ShipmentLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, kernel.ShipmentLine{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	shipment, err := s.ship.CreateShipment(kernel.CreateShipmentRecord{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Lines:      lines,
	})
	if err != nil {
		return kernel.CreateShipmentResult{}, err
	}

	if err := order.MarkShipped(s.clock.Now()); err != nil {
		return kernel.CreateShipmentResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return kernel.CreateShipmentResult{}, err
	}

	return kernel.CreateShipmentResult{
		ShipmentID: shipment.ShipmentID,
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) CancelOrder(command kernel.CancelOrderCommand) (kernel.CancelOrderResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return kernel.CancelOrderResult{}, err
	}

	items := make([]kernel.InventoryReservationItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		items = append(items, kernel.InventoryReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := order.Cancel(); err != nil {
		return kernel.CancelOrderResult{}, err
	}

	if err := s.release.Release(items); err != nil {
		return kernel.CancelOrderResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return kernel.CancelOrderResult{}, err
	}

	return kernel.CancelOrderResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) GetReturnableOrder(orderID string) (kernel.ReturnableOrder, error) {
	order, err := s.orders.FindByID(orderID)
	if err != nil {
		return kernel.ReturnableOrder{}, err
	}

	if err := order.EnsureReturnable(); err != nil {
		return kernel.ReturnableOrder{}, err
	}

	lines := make([]kernel.ReturnableOrderLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, kernel.ReturnableOrderLine{
			ProductSKU:       line.ProductSKU,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return kernel.ReturnableOrder{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		ShippedAt:  order.ShippedAt,
		Lines:      lines,
	}, nil
}

func (s Service) GetOrder(query kernel.GetOrderQuery) (kernel.OrderDetails, error) {
	order, err := s.orders.FindByID(query.OrderID)
	if err != nil {
		return kernel.OrderDetails{}, err
	}

	lines := make([]kernel.OrderLineDetails, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, kernel.OrderLineDetails{
			ProductSKU:      line.ProductSKU,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
		})
	}

	return kernel.OrderDetails{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
		Lines:      lines,
	}, nil
}

func (s Service) ListOrders(query kernel.ListOrdersQuery) ([]kernel.OrderSummary, error) {
	orders, err := s.orders.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	results := make([]kernel.OrderSummary, 0, len(orders))
	for _, order := range orders {
		results = append(results, kernel.OrderSummary{
			OrderID:    order.ID,
			QuoteID:    order.QuoteID,
			CustomerID: order.CustomerID,
			Status:     order.Status,
			LineCount:  len(order.Lines),
		})
	}

	return results, nil
}
