package main

import (
	"fmt"
	"log"

	"rich-domain-model-architecture/internal/domain/catalog"
	"rich-domain-model-architecture/internal/domain/customer"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func main() {
	customerAggregate, err := customer.NewCustomer("customer-001", customer.CustomerTierPreferred, customer.PaymentTermsInvoice30)
	if err != nil {
		log.Fatal(err)
	}
	if err := customerAggregate.EnsureCanCreateQuote(); err != nil {
		log.Fatal(err)
	}

	quote, err := quoting.NewQuote("quote-001", quoting.CustomerID(customerAggregate.ID()))
	if err != nil {
		log.Fatal(err)
	}

	catalogPrice, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Desk", catalog.ProductCategoryStandard, catalogPrice)
	if err != nil {
		log.Fatal(err)
	}
	if err := product.EnsureSellable(); err != nil {
		log.Fatal(err)
	}

	quotePrice, err := quoting.NewMoney(product.BasePrice().Cents(), product.BasePrice().Currency())
	if err != nil {
		log.Fatal(err)
	}
	line, err := quoting.NewQuoteLineFromProductSnapshot(
		quoting.ProductSKU(product.SKU()),
		product.Name(),
		2,
		quotePrice,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		log.Fatal(err)
	}

	newCatalogPrice, err := catalog.NewPrice(18000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	if err := product.ChangePrice(newCatalogPrice); err != nil {
		log.Fatal(err)
	}

	total, err := quote.Total()
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.SubmitForApproval(); err != nil {
		log.Fatal(err)
	}

	snapshot := quote.Lines()[0]
	fmt.Printf("customer domain object: id=%s tier=%s terms=%s active=%t\n", customerAggregate.ID(), customerAggregate.Tier(), customerAggregate.PaymentTerms(), customerAggregate.Active())
	fmt.Printf("product domain object: sku=%s active=%t price=%d %s\n", product.SKU(), product.Active(), product.BasePrice().Cents(), product.BasePrice().Currency())
	fmt.Printf("quote snapshot: sku=%s name=%s unit-price=%d %s\n", snapshot.ProductSKU(), snapshot.ProductName(), snapshot.UnitPrice().Cents(), snapshot.UnitPrice().Currency())
	fmt.Printf("domain aggregate: id=%s customer=%s status=%s lines=%d total=%d %s\n", quote.ID(), quote.CustomerID(), quote.Status(), len(quote.Lines()), total.Cents(), total.Currency())

	if err := quote.Approve(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("approved aggregate: id=%s status=%s\n", quote.ID(), quote.Status())
}
