package clock

import "microkernel-architecture/internal/kernel"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

func (p Plugin) ID() string {
	return "clock"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposeClock(NewService())
	return nil
}
