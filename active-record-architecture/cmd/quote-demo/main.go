package main

import (
	"fmt"
	"log"

	"active-record-architecture/internal/records"
)

func main() {
	db := records.NewDatabase()

	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		log.Fatal(err)
	}

	loadedCustomer, err := records.FindCustomer(db, customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	quote, err := records.NewDraftQuote(db, loadedCustomer.ID)
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.Save(); err != nil {
		log.Fatal(err)
	}

	reloadedQuote, err := records.FindQuote(db, quote.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"created draft quote: id=%s customer=%s status=%s\n",
		reloadedQuote.ID,
		reloadedQuote.CustomerID,
		reloadedQuote.Status,
	)
}
