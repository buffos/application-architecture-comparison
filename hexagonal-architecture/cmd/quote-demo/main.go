package main

import (
	"fmt"
	"log"

	cli "hexagonal-architecture/internal/adapters/cli"
	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/adapters/services/approval"
	"hexagonal-architecture/internal/adapters/services/pricing"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func main() {
	quoteRepo := memory.NewQuoteRepository()
	orderRepo := memory.NewOrderRepository()
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	inventory := memory.NewInventoryReservationAdapter(map[string]int{
		"CHAIR-001": 5,
	})
	pricingPolicy := pricing.NewFixedPricingPolicy()
	approvalPolicy := approval.NewCategoryApprovalPolicy()
	if err := customerRepo.Save(domain.Customer{ID: "customer-001", Active: true}); err != nil {
		log.Fatal(err)
	}
	if err := productRepo.Save(domain.Product{
		SKU:       "CHAIR-001",
		Name:      "Office Chair",
		Category:  "Standard",
		BasePrice: 10000,
		Available: true,
	}); err != nil {
		log.Fatal(err)
	}

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	addQuoteLine := application.NewAddQuoteLineUseCase(quoteRepo, productRepo, pricingPolicy)
	submitQuote := application.NewSubmitQuoteUseCase(quoteRepo, approvalPolicy)
	convertQuote := application.NewConvertQuoteToOrderUseCase(quoteRepo, orderRepo, inventory)
	getQuote := application.NewGetQuoteUseCase(quoteRepo)
	handler := cli.NewQuoteHandler(createQuote, addQuoteLine, submitQuote, convertQuote, getQuote)

	output, err := handler.RunDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
