package main

import (
	"fmt"
	"log"

	"onion-architecture/internal/application"
	"onion-architecture/internal/domain"
	"onion-architecture/internal/infrastructure/memory"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	quoteRepository := memory.NewQuoteRepository()

	if err := customerRepository.Save(domain.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	service := application.NewCreateDraftQuoteService(quoteRepository, customerRepository)
	getQuote := application.NewGetQuoteService(quoteRepository)

	result, err := service.Execute(application.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	details, err := getQuote.Execute(application.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s\n", details.QuoteID, details.CustomerID, details.Status)
}
