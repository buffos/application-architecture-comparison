package payments

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	gateway Gateway
}

func NewPlugin(gateway Gateway) Plugin {
	return Plugin{gateway: gateway}
}

func (p Plugin) ID() string {
	return "payments"
}

func (p Plugin) Register(host *kernel.Host) error {
	service := NewService(p.gateway)
	host.ExposePaymentCapture(service)
	host.ExposePaymentRefund(service)
	return nil
}
