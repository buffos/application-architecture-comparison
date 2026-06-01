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
	createPresenter := presenters.NewCreateDraftQuotePresenter()
	getPresenter := presenters.NewGetQuotePresenter()

	if err := customerGateway.Save(entities.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	createInteractor := usecases.NewCreateDraftQuoteInteractor(quoteGateway, customerGateway, createPresenter)
	createController := controllers.NewCreateDraftQuoteController(createInteractor)

	if err := createController.Handle("customer-001"); err != nil {
		log.Fatal(err)
	}

	getInteractor := usecases.NewGetQuoteInteractor(quoteGateway, getPresenter)
	getController := controllers.NewGetQuoteController(getInteractor)

	if err := getController.Handle(createPresenter.ViewModel().QuoteID); err != nil {
		log.Fatal(err)
	}

	fmt.Println(createPresenter.ViewModel().Message)
	fmt.Println(getPresenter.ViewModel().Message)
}
