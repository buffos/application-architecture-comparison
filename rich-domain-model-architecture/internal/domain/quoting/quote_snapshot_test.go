package quoting_test

import (
	"testing"

	"rich-domain-model-architecture/internal/domain/catalog"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestQuoteLineKeepsAProductSnapshot(t *testing.T) {
	catalogPrice, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Desk", catalog.ProductCategoryStandard, catalogPrice)
	if err != nil {
		t.Fatal(err)
	}

	quotePrice, err := quoting.NewMoney(product.BasePrice().Cents(), product.BasePrice().Currency())
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLineFromProductSnapshot(
		quoting.ProductSKU(product.SKU()),
		product.Name(),
		2,
		quotePrice,
	)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}

	newCatalogPrice, err := catalog.NewPrice(18000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if err := product.ChangePrice(newCatalogPrice); err != nil {
		t.Fatal(err)
	}
	if err := product.Discontinue(); err != nil {
		t.Fatal(err)
	}

	storedLine := quote.Lines()[0]
	if storedLine.ProductSKU() != "sku-001" {
		t.Fatalf("stored sku = %s, want sku-001", storedLine.ProductSKU())
	}
	if storedLine.ProductName() != "Desk" {
		t.Fatalf("stored name = %s, want Desk", storedLine.ProductName())
	}
	if storedLine.UnitPrice().Cents() != 15000 {
		t.Fatalf("stored unit price = %d, want 15000", storedLine.UnitPrice().Cents())
	}
}
