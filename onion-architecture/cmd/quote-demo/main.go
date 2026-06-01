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
	productRepository := memory.NewProductRepository()

	if err := customerRepository.Save(domain.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(domain.Product{
		SKU:      "sku-001",
		Name:     "Desk",
		Active:   true,
		UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	service := application.NewCreateDraftQuoteService(quoteRepository, customerRepository)
	getQuote := application.NewGetQuoteService(quoteRepository)
	addQuoteLine := application.NewAddQuoteLineService(quoteRepository, productRepository)

	result, err := service.Execute(application.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := addQuoteLine.Execute(application.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	details, err := getQuote.Execute(application.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)
}
