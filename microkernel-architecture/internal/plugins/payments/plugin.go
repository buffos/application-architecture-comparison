package payments

import "microkernel-architecture/internal/kernel"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

func (p Plugin) ID() string {
	return "payments"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposePaymentCapture(NewService())
	return nil
}
