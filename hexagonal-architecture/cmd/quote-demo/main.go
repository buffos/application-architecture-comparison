package main

import (
	"fmt"
	"log"

	cli "hexagonal-architecture/internal/adapters/cli"
	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/application"
)

func main() {
	repo := memory.NewQuoteRepository()
	createQuote := application.NewCreateDraftQuoteUseCase(repo)
	getQuote := application.NewGetQuoteUseCase(repo)
	handler := cli.NewQuoteHandler(createQuote, getQuote)

	output, err := handler.RunDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
