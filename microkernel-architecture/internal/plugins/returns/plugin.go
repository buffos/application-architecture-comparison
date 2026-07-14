package returns

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	requests Repository
}

func NewPlugin(requests Repository) Plugin {
	return Plugin{
		requests: requests,
	}
}

func (p Plugin) ID() string {
	return "returns"
}

func (p Plugin) Register(host *kernel.Host) error {
	orders, err := host.ReturnableOrderProvider()
	if err != nil {
		return err
	}

	refunds, err := host.PaymentRefund()
	if err != nil {
		return err
	}

	host.ExposeReturnService(NewService(p.requests, orders, refunds))
	return nil
}
