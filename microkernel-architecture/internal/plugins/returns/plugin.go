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

	clock, err := host.Clock()
	if err != nil {
		return err
	}

	refunds, err := host.PaymentRefund()
	if err != nil {
		return err
	}

	policy, err := host.ReturnEligibilityPolicy()
	if err != nil {
		return err
	}

	idempotency, err := host.IdempotencyStore()
	if err != nil {
		return err
	}

	restock, err := host.InventoryRestock()
	if err != nil {
		return err
	}

	host.ExposeReturnService(NewService(p.requests, orders, clock, policy, idempotency, refunds, restock))
	return nil
}
