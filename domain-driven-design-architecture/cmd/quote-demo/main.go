package main

import (
	"fmt"
	"log"

	"domain-driven-design-architecture/internal/domain/quoting"
)

func main() {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		log.Fatal(err)
	}
	unitPrice, err := quoting.NewMoney(15000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, unitPrice)
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
	fmt.Printf("quote aggregate: id=%s status=%s lines=%d total=%d %s\n", quote.ID(), quote.Status(), len(quote.Lines()), total.Cents(), total.Currency())
	if err := quote.Submit(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("submitted quote aggregate: id=%s status=%s\n", quote.ID(), quote.Status())
}
