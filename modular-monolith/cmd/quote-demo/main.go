package main

import (
	"fmt"
	"log"

	"modular-monolith/internal/modules/customers"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/platform/memory"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	quoteRepository := memory.NewQuoteRepository()

	if err := customerRepository.Save(customers.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	customerModule := customers.NewService(customerRepository)
	quoteModule := quotes.NewService(quoteRepository, customerModule)

	result, err := quoteModule.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	details, err := quoteModule.GetQuote(quotes.GetQuoteQuery{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s\n", details.QuoteID, details.CustomerID, details.Status)
}
