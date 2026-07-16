package main

import (
	"fmt"
	"log"

	"component-based-architecture/internal/components/customers"
	"component-based-architecture/internal/components/products"
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
	productComponent := products.NewComponent()
	if err := productComponent.Register(products.Product{
		SKU: "sku-001", Name: "Desk", Active: true, UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	quoteComponent := quotes.NewComponent(customerComponent, productComponent)
	result, err := quoteComponent.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := quoteComponent.AddQuoteLine(quotes.AddQuoteLineCommand{
		QuoteID: result.QuoteID, ProductSKU: "sku-001", Quantity: 2,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("added quote line: id=%s lines=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.Status)

	submission, err := quoteComponent.SubmitQuote(quotes.SubmitQuoteCommand{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("submitted quote: id=%s lines=%d status=%s\n", submission.QuoteID, submission.LineCount, submission.Status)

	var quoteLookup quotes.QuoteLookup = quoteComponent
	details, err := quoteLookup.GetQuote(quotes.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)
}
