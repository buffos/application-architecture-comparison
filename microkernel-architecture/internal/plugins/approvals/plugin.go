package approvals

import "microkernel-architecture/internal/kernel"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

func (p Plugin) ID() string {
	return "approvals"
}

func (p Plugin) Register(host *kernel.Host) error {
	host.ExposeApprovalPolicy(NewService())
	return nil
}
