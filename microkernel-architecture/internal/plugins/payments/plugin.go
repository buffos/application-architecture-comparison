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
	service := NewService()
	host.ExposePaymentCapture(service)
	host.ExposePaymentRefund(service)
	return nil
}
