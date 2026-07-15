package customers

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	customers Repository
}

func NewPlugin(customers Repository) Plugin {
	return Plugin{
		customers: customers,
	}
}

func (p Plugin) ID() string {
	return "customers"
}

func (p Plugin) Register(host *kernel.Host) error {
	service := NewService(p.customers)
	host.ExposeCustomerDirectory(service)
	host.ExposeCustomerReader(service)
	return nil
}
