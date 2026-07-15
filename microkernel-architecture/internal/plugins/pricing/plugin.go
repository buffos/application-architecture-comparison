package pricing

import "microkernel-architecture/internal/kernel"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

func (p Plugin) ID() string {
	return "pricing"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposeQuotePricer(NewService())
	return nil
}
