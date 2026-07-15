package reporting

import "microkernel-architecture/internal/kernel"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

func (p Plugin) ID() string {
	return "reporting"
}

func (p Plugin) Register(host *kernel.Host) error {
	quotes, err := host.QuoteReader()
	if err != nil {
		return err
	}

	orders, err := host.OrderReader()
	if err != nil {
		return err
	}

	returns, err := host.ReturnReader()
	if err != nil {
		return err
	}

	host.ExposeReporting(NewService(quotes, orders, returns))
	return nil
}
