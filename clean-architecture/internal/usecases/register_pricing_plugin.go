package usecases

import "clean-architecture/internal/entities"

type RegisterPricingPluginInput struct {
	Name string
}

type RegisterPricingPluginOutput struct {
	Name    string
	Enabled bool
}

type RegisterPricingPluginInputBoundary interface {
	Execute(input RegisterPricingPluginInput) error
}

type RegisterPricingPluginOutputBoundary interface {
	Present(output RegisterPricingPluginOutput) error
}

type PluginWriter interface {
	Save(plugin entities.PluginRegistration) error
}

type RegisterPricingPluginInteractor struct {
	plugins PluginWriter
	output  RegisterPricingPluginOutputBoundary
}

func NewRegisterPricingPluginInteractor(plugins PluginWriter, output RegisterPricingPluginOutputBoundary) RegisterPricingPluginInteractor {
	return RegisterPricingPluginInteractor{plugins: plugins, output: output}
}

func (uc RegisterPricingPluginInteractor) Execute(input RegisterPricingPluginInput) error {
	plugin, err := entities.NewPluginRegistration(input.Name)
	if err != nil {
		return err
	}

	if err := uc.plugins.Save(plugin); err != nil {
		return err
	}

	return uc.output.Present(RegisterPricingPluginOutput{
		Name:    plugin.Name,
		Enabled: plugin.Enabled,
	})
}
