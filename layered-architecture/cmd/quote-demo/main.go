package main

import (
	"fmt"
	"log"

	"layered-architecture/internal/application"
	"layered-architecture/internal/infrastructure/memory"
)

func main() {
	repo := memory.NewQuoteRepository()
	service := application.NewQuoteService(repo)

	quote, err := service.CreateDraftQuote("customer-001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", quote.ID, quote.CustomerID, quote.Status)
}
