package idempotency

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	store Store
}

func NewPlugin(store Store) Plugin {
	return Plugin{
		store: store,
	}
}

func (p Plugin) ID() string {
	return "idempotency"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposeIdempotencyStore(NewService(p.store))
	return nil
}
