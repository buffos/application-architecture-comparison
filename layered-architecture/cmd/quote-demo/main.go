package main

import (
	"fmt"
	"log"

	"layered-architecture/internal/application"
	"layered-architecture/internal/infrastructure/memory"
	"layered-architecture/internal/presentation/console"
)

func main() {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()

	customerService := application.NewCustomerService(customerRepo)
	catalogService := application.NewCatalogService(productRepo)
	quoteService := application.NewQuoteService(quoteRepo, customerRepo, productRepo)
	handler := console.NewQuoteHandler(customerService, catalogService, quoteService)

	output, err := handler.RunDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
