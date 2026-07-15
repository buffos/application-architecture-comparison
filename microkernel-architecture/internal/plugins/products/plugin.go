package products

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	products Repository
}

func NewPlugin(products Repository) Plugin {
	return Plugin{
		products: products,
	}
}

func (p Plugin) ID() string {
	return "products"
}

func (p Plugin) Register(host *kernel.Host) error {
	service := NewService(p.products)
	host.ExposeProductCatalog(service)
	host.ExposeProductReader(service)
	return nil
}
