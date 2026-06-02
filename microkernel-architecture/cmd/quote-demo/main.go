package main

import (
	"fmt"
	"log"

	"microkernel-architecture/internal/kernel"
	"microkernel-architecture/internal/platform/memory"
	"microkernel-architecture/internal/plugins/customers"
	"microkernel-architecture/internal/plugins/products"
	"microkernel-architecture/internal/plugins/quotes"
)

func main() {
	host := kernel.NewHost()

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
		Active:    true,
		UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(customers.NewPlugin(customerRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(products.NewPlugin(productRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(quotes.NewPlugin(quoteRepository)); err != nil {
		log.Fatal(err)
	}

	quoteService, err := host.QuoteService()
	if err != nil {
		log.Fatal(err)
	}

	quoteReader, err := host.QuoteReader()
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

	lineResult, err := quoteService.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	details, err := quoteReader.GetQuote(kernel.GetQuoteQuery{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d items=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount, details.TotalItems)
}
