package main

import (
	"fmt"
	"log"

	"domain-driven-design-architecture/internal/domain/catalog"
	"domain-driven-design-architecture/internal/domain/customer"
	"domain-driven-design-architecture/internal/domain/inventory"
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
	productPrice, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Desk", catalog.ProductCategoryStandard, productPrice, 30)
	if err != nil {
		log.Fatal(err)
	}
	if err := product.EnsureSellable(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("product aggregate: sku=%s category=%s active=%t\n", product.SKU(), product.Category(), product.Active())
	stock, err := inventory.NewStockRecord(inventory.ProductSKU(product.SKU()), 10, 3)
	if err != nil {
		log.Fatal(err)
	}
	if err := stock.Reserve(2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("stock aggregate: sku=%s available=%d reserved=%d low=%t\n", stock.SKU(), stock.Available(), stock.Reserved(), stock.IsLowStock())
	quotePrice, err := quoting.NewMoney(product.BasePrice().Cents(), product.BasePrice().Currency())
	if err != nil {
		log.Fatal(err)
	}
	line, err := quoting.NewQuoteLine(quoting.ProductSKU(product.SKU()), 2, quotePrice)
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
