package main

import (
	"fmt"
	"log"

	"modular-monolith/internal/modules/approvals"
	"modular-monolith/internal/modules/customers"
	"modular-monolith/internal/modules/products"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/platform/memory"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	productRepository := memory.NewProductRepository()
	quoteRepository := memory.NewQuoteRepository()

	if err := customerRepository.Save(customers.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(products.Product{
		SKU:       "sku-001",
		Name:      "Desk",
		Category:  "Standard",
		Active:    true,
		UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	customerModule := customers.NewService(customerRepository)
	productModule := products.NewService(productRepository)
	approvalModule := approvals.NewService()
	quoteModule := quotes.NewService(quoteRepository, customerModule, productModule, approvalModule)

	result, err := quoteModule.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := quoteModule.AddQuoteLine(quotes.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	submitResult, err := quoteModule.SubmitQuote(quotes.SubmitQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted quote: id=%s lines=%d items=%d status=%s\n", submitResult.QuoteID, submitResult.LineCount, submitResult.TotalItems, submitResult.Status)

	details, err := quoteModule.GetQuote(quotes.GetQuoteQuery{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)
}
