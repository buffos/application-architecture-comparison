package main

import (
	"fmt"
	"log"

	"onion-architecture/internal/application"
	"onion-architecture/internal/domain"
	"onion-architecture/internal/infrastructure/memory"
	"onion-architecture/internal/infrastructure/policies/approval"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	quoteRepository := memory.NewQuoteRepository()
	productRepository := memory.NewProductRepository()
	orderRepository := memory.NewOrderRepository()

	if err := customerRepository.Save(domain.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(domain.Product{
		SKU:       "sku-001",
		Name:      "Desk",
		Category:  "Standard",
		Active:    true,
		UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(domain.Product{
		SKU:       "sku-002",
		Name:      "Custom Desk",
		Category:  "CustomBuild",
		Active:    true,
		UnitPrice: 45000,
	}); err != nil {
		log.Fatal(err)
	}

	submissionPolicy := approval.NewCategoryPolicy()
	service := application.NewCreateDraftQuoteService(quoteRepository, customerRepository)
	getQuote := application.NewGetQuoteService(quoteRepository)
	addQuoteLine := application.NewAddQuoteLineService(quoteRepository, productRepository)
	submitQuote := application.NewSubmitQuoteService(quoteRepository, submissionPolicy)
	approveQuote := application.NewApproveQuoteService(quoteRepository)
	convertQuote := application.NewConvertQuoteToOrderService(quoteRepository, orderRepository)

	result, err := service.Execute(application.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := addQuoteLine.Execute(application.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-002",
		Quantity:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	submitResult, err := submitQuote.Execute(application.SubmitQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted quote: id=%s lines=%d items=%d status=%s\n", submitResult.QuoteID, submitResult.LineCount, submitResult.TotalItems, submitResult.Status)

	approvalResult, err := approveQuote.Execute(application.ApproveQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("approved quote: id=%s lines=%d items=%d status=%s\n", approvalResult.QuoteID, approvalResult.LineCount, approvalResult.TotalItems, approvalResult.Status)

	orderResult, err := convertQuote.Execute(application.ConvertQuoteToOrderCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("converted quote: order=%s quote=%s customer=%s status=%s lines=%d\n", orderResult.OrderID, orderResult.QuoteID, orderResult.CustomerID, orderResult.Status, orderResult.LineCount)

	details, err := getQuote.Execute(application.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)
}
