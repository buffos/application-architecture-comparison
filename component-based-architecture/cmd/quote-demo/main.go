package main

import (
	"fmt"
	"log"

	"component-based-architecture/internal/components/customers"
	"component-based-architecture/internal/components/quotes"
)

func main() {
	customerComponent := customers.NewComponent()
	if err := customerComponent.Register(customers.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	quoteComponent := quotes.NewComponent(customerComponent)
	result, err := quoteComponent.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	var quoteLookup quotes.QuoteLookup = quoteComponent
	details, err := quoteLookup.GetQuote(quotes.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s\n", details.QuoteID, details.CustomerID, details.Status)
}
