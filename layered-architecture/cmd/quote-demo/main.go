package main

import (
	"fmt"
	"log"

	"layered-architecture/internal/application"
	"layered-architecture/internal/infrastructure/memory"
	"layered-architecture/internal/presentation/console"
)

func main() {
	repo := memory.NewQuoteRepository()
	service := application.NewQuoteService(repo)
	handler := console.NewQuoteHandler(service)

	output, err := handler.RunDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
