package quotes

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	quotes Repository
}

func NewPlugin(quotes Repository) Plugin {
	return Plugin{
		quotes: quotes,
	}
}

func (p Plugin) ID() string {
	return "quotes"
}

func (p Plugin) Register(host *kernel.Host) error {
	customers, err := host.CustomerDirectory()
	if err != nil {
		return err
	}

	service := NewService(p.quotes, customers)
	host.ExposeQuoteService(service)
	host.ExposeQuoteReader(service)
	return nil
}
