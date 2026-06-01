package controllers

import "clean-architecture/internal/usecases"

type ListPluginsController struct {
	useCase usecases.ListPluginsInputBoundary
}

func NewListPluginsController(useCase usecases.ListPluginsInputBoundary) ListPluginsController {
	return ListPluginsController{useCase: useCase}
}

func (c ListPluginsController) Handle() error {
	return c.useCase.Execute(usecases.ListPluginsInput{})
}
