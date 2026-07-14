package orders

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	orders Repository
}

func NewPlugin(orders Repository) Plugin {
	return Plugin{
		orders: orders,
	}
}

func (p Plugin) ID() string {
	return "orders"
}

func (p Plugin) Register(host *kernel.Host) error {
	quotes, err := host.ApprovedQuoteProvider()
	if err != nil {
		return err
	}

	stock, err := host.InventoryReservation()
	if err != nil {
		return err
	}

	payments, err := host.PaymentCapture()
	if err != nil {
		return err
	}

	host.ExposeOrderService(NewService(p.orders, quotes, stock, payments))
	return nil
}
