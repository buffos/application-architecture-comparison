package main

import (
	"fmt"
	"log"

	"clean-architecture/internal/entities"
	"clean-architecture/internal/infrastructure/memory"
	"clean-architecture/internal/interfaceadapters/controllers"
	"clean-architecture/internal/interfaceadapters/presenters"
	"clean-architecture/internal/usecases"
)

func main() {
	customerGateway := memory.NewCustomerGateway()
	quoteGateway := memory.NewQuoteGateway()
	presenter := presenters.NewCreateDraftQuotePresenter()

	if err := customerGateway.Save(entities.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	interactor := usecases.NewCreateDraftQuoteInteractor(quoteGateway, customerGateway, presenter)
	controller := controllers.NewCreateDraftQuoteController(interactor)

	if err := controller.Handle("customer-001"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(presenter.ViewModel().Message)
}
