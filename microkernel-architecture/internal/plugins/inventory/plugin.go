package inventory

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	stock Repository
}

func NewPlugin(stock Repository) Plugin {
	return Plugin{
		stock: stock,
	}
}

func (p Plugin) ID() string {
	return "inventory"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposeInventoryReservation(NewService(p.stock))
	return nil
}
