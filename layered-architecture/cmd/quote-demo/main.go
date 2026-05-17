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

	createdQuote, err := service.CreateDraftQuote("customer-001")
	if err != nil {
		log.Fatal(err)
	}

	quoteWithLine, err := service.AddQuoteLine(createdQuote.ID, "Office Chair", 2)
	if err != nil {
		log.Fatal(err)
	}

	submittedQuote, err := service.SubmitQuote(createdQuote.ID)
	if err != nil {
		log.Fatal(err)
	}

	loadedQuote, err := service.GetQuote(createdQuote.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", createdQuote.ID, createdQuote.CustomerID, createdQuote.Status)
	fmt.Printf("added quote line: id=%s lines=%d status=%s\n", quoteWithLine.ID, len(quoteWithLine.Lines), quoteWithLine.Status)
	fmt.Printf("submitted quote: id=%s lines=%d status=%s\n", submittedQuote.ID, len(submittedQuote.Lines), submittedQuote.Status)
	fmt.Printf("loaded draft quote: id=%s customer=%s status=%s\n", loadedQuote.ID, loadedQuote.CustomerID, loadedQuote.Status)
}
