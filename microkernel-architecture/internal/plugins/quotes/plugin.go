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

	host.ExposeQuoteService(NewService(p.quotes, customers))
	return nil
}
