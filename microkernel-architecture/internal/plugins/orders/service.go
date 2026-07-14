package orders

import "microkernel-architecture/internal/kernel"

type Service struct {
	orders Repository
	quotes kernel.ApprovedQuoteProvider
	stock  kernel.InventoryReservation
	pay    kernel.PaymentCapture
	ship   kernel.ShipmentCreation
}

func NewService(orders Repository, quotes kernel.ApprovedQuoteProvider, stock kernel.InventoryReservation, pay kernel.PaymentCapture, ship kernel.ShipmentCreation) Service {
	return Service{
		orders: orders,
		quotes: quotes,
		stock:  stock,
		pay:    pay,
		ship:   ship,
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
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
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

	if err := order.MarkShipped(); err != nil {
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
