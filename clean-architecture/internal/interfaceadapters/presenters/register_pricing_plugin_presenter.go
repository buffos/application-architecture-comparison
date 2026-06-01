package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type RegisterPricingPluginViewModel struct {
	Message string
	Name    string
	Enabled bool
}

type RegisterPricingPluginPresenter struct {
	viewModel RegisterPricingPluginViewModel
}

func NewRegisterPricingPluginPresenter() *RegisterPricingPluginPresenter {
	return &RegisterPricingPluginPresenter{}
}

func (p *RegisterPricingPluginPresenter) Present(output usecases.RegisterPricingPluginOutput) error {
	p.viewModel = RegisterPricingPluginViewModel{
		Message: fmt.Sprintf("registered pricing plugin: name=%s enabled=%t", output.Name, output.Enabled),
		Name:    output.Name,
		Enabled: output.Enabled,
	}
	return nil
}

func (p *RegisterPricingPluginPresenter) ViewModel() RegisterPricingPluginViewModel {
	return p.viewModel
}
