package main

import (
	"fmt"
	"log"

	"microkernel-architecture/internal/kernel"
	"microkernel-architecture/internal/platform/memory"
	"microkernel-architecture/internal/plugins/customers"
	"microkernel-architecture/internal/plugins/quotes"
)

func main() {
	host := kernel.NewHost()

	customerRepository := memory.NewCustomerRepository()
	quoteRepository := memory.NewQuoteRepository()

	if err := customerRepository.Save(customers.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(customers.NewPlugin(customerRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(quotes.NewPlugin(quoteRepository)); err != nil {
		log.Fatal(err)
	}

	quoteService, err := host.QuoteService()
	if err != nil {
		log.Fatal(err)
	}

	result, err := quoteService.CreateDraftQuote(kernel.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)
}
