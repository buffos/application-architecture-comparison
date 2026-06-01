package usecases

import "clean-architecture/internal/entities"

type EnablePluginInput struct {
	Name string
}

type EnablePluginOutput struct {
	Name    string
	Enabled bool
}

type EnablePluginInputBoundary interface {
	Execute(input EnablePluginInput) error
}

type EnablePluginOutputBoundary interface {
	Present(output EnablePluginOutput) error
}

type PluginEditor interface {
	FindByName(name string) (entities.PluginRegistration, error)
	Save(plugin entities.PluginRegistration) error
}

type EnablePluginInteractor struct {
	plugins PluginEditor
	output  EnablePluginOutputBoundary
}

func NewEnablePluginInteractor(plugins PluginEditor, output EnablePluginOutputBoundary) EnablePluginInteractor {
	return EnablePluginInteractor{plugins: plugins, output: output}
}

func (uc EnablePluginInteractor) Execute(input EnablePluginInput) error {
	plugin, err := uc.plugins.FindByName(input.Name)
	if err != nil {
		return err
	}

	plugin.Enable()
	if err := uc.plugins.Save(plugin); err != nil {
		return err
	}

	return uc.output.Present(EnablePluginOutput{Name: plugin.Name, Enabled: plugin.Enabled})
}
