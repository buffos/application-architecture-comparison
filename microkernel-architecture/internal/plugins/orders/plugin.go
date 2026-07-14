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

	release, err := host.InventoryRelease()
	if err != nil {
		return err
	}

	payments, err := host.PaymentCapture()
	if err != nil {
		return err
	}

	shipments, err := host.ShipmentCreation()
	if err != nil {
		return err
	}

	clock, err := host.Clock()
	if err != nil {
		return err
	}

	service := NewService(p.orders, quotes, stock, release, payments, shipments, clock)
	host.ExposeOrderService(service)
	host.ExposeOrderReader(service)
	host.ExposeReturnableOrderProvider(service)
	return nil
}
