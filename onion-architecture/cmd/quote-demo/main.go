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

	result, err := service.Execute(application.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)
}
