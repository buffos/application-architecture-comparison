package returneligibility

import "microkernel-architecture/internal/kernel"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

func (p Plugin) ID() string {
	return "returneligibility"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposeReturnEligibilityPolicy(NewService())
	return nil
}
