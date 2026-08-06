package main

import (
	"fmt"
	"log"

	"transaction-script-architecture/internal/data"
	"transaction-script-architecture/internal/scripts"
)

func main() {
	store := data.NewStore()
	store.Customers["customer-001"] = data.Customer{
		ID:     "customer-001",
		Active: true,
	}

	quote, err := scripts.CreateDraftQuote(store, "customer-001")
	if err != nil {
		log.Fatal(err)
	}

	_, saved := store.Quotes[quote.ID]
	fmt.Printf(
		"created draft quote: id=%s customer=%s status=%s saved=%t\n",
		quote.ID,
		quote.CustomerID,
		quote.Status,
		saved,
	)
}
