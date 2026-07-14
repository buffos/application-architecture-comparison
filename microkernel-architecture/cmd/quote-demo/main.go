package main

import (
	"fmt"
	"log"

	"microkernel-architecture/internal/kernel"
	"microkernel-architecture/internal/platform/memory"
	"microkernel-architecture/internal/plugins/approvals"
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
		Category:  "Standard",
		Active:    true,
		UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(products.Product{
		SKU:       "sku-002",
		Name:      "Custom Desk",
		Category:  "CustomBuild",
		Active:    true,
		UnitPrice: 45000,
	}); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(customers.NewPlugin(customerRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(products.NewPlugin(productRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(approvals.NewPlugin()); err != nil {
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

	submitResult, err := quoteService.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted quote: id=%s lines=%d items=%d status=%s\n", submitResult.QuoteID, submitResult.LineCount, submitResult.TotalItems, submitResult.Status)

	details, err := quoteReader.GetQuote(kernel.GetQuoteQuery{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d items=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount, details.TotalItems)

	pendingResult, err := quoteService.CreateDraftQuote(kernel.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = quoteService.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    pendingResult.QuoteID,
		ProductSKU: "sku-002",
		Quantity:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	pendingSubmit, err := quoteService.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted custom quote: id=%s status=%s\n", pendingSubmit.QuoteID, pendingSubmit.Status)

	approvedPending, err := quoteService.ApproveQuote(kernel.ApproveQuoteCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("approved pending quote: id=%s status=%s\n", approvedPending.QuoteID, approvedPending.Status)
}
