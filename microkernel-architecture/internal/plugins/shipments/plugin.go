package shipments

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	shipments Repository
}

func NewPlugin(shipments Repository) Plugin {
	return Plugin{
		shipments: shipments,
	}
}

func (p Plugin) ID() string {
	return "shipments"
}

func (p Plugin) Register(host *kernel.Host) error {
	service := NewService(p.shipments)
	host.ExposeShipmentCreation(service)
	host.ExposeShipmentReader(service)
	return nil
}
