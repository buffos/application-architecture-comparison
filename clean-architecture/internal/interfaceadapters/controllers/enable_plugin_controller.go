package controllers

import "clean-architecture/internal/usecases"

type EnablePluginController struct {
	useCase usecases.EnablePluginInputBoundary
}

func NewEnablePluginController(useCase usecases.EnablePluginInputBoundary) EnablePluginController {
	return EnablePluginController{useCase: useCase}
}

func (c EnablePluginController) Handle(name string) error {
	return c.useCase.Execute(usecases.EnablePluginInput{Name: name})
}
