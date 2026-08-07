package main

import (
	"fmt"
	"log"

	"rich-domain-model-architecture/internal/domain/quoting"
)

func main() {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		log.Fatal(err)
	}

	price, err := quoting.NewMoney(15000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		log.Fatal(err)
	}

	total, err := quote.Total()
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.SubmitForApproval(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("domain aggregate: id=%s customer=%s status=%s lines=%d total=%d %s\n", quote.ID(), quote.CustomerID(), quote.Status(), len(quote.Lines()), total.Cents(), total.Currency())

	if err := quote.Approve(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("approved aggregate: id=%s status=%s\n", quote.ID(), quote.Status())
}
