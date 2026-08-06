package main

import (
	"fmt"
	"log"

	"domain-driven-design-architecture/internal/domain/customer"
	"domain-driven-design-architecture/internal/domain/quoting"
)

func main() {
	customerAggregate, err := customer.NewCustomer("customer-001", customer.CustomerTierPreferred, customer.PaymentTermsInvoice30)
	if err != nil {
		log.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-001", quoting.CustomerID(customerAggregate.ID()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("customer aggregate: id=%s tier=%s active=%t\n", customerAggregate.ID(), customerAggregate.Tier(), customerAggregate.Active())
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
