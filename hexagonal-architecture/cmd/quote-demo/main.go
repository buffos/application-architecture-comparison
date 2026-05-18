package main

import (
	"fmt"
	"log"

	cli "hexagonal-architecture/internal/adapters/cli"
	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

func main() {
	quoteRepo := memory.NewQuoteRepository()
	customerRepo := memory.NewCustomerRepository()
	if err := customerRepo.Save(domain.Customer{ID: "customer-001", Active: true}); err != nil {
		log.Fatal(err)
	}

	createQuote := application.NewCreateDraftQuoteUseCase(quoteRepo, customerRepo)
	getQuote := application.NewGetQuoteUseCase(quoteRepo)
	handler := cli.NewQuoteHandler(createQuote, getQuote)

	output, err := handler.RunDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
