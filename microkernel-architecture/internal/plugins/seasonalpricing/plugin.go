package seasonalpricing

import "microkernel-architecture/internal/kernel"

type Plugin struct {
	discountPercent int
}

func NewPlugin(discountPercent int) Plugin {
	return Plugin{discountPercent: discountPercent}
}

func (p Plugin) ID() string {
	return "seasonalpricing"
}

func (p Plugin) Register(host *kernel.Host) error {
	base, err := host.QuotePricer()
	if err != nil {
		return err
	}

	host.ExposeQuotePricer(NewService(base, p.discountPercent))
	return nil
}
