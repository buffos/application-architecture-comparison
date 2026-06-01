package controllers

import "clean-architecture/internal/usecases"

type RegisterPricingPluginController struct {
	useCase usecases.RegisterPricingPluginInputBoundary
}

func NewRegisterPricingPluginController(useCase usecases.RegisterPricingPluginInputBoundary) RegisterPricingPluginController {
	return RegisterPricingPluginController{useCase: useCase}
}

func (c RegisterPricingPluginController) Handle(name string) error {
	return c.useCase.Execute(usecases.RegisterPricingPluginInput{Name: name})
}
