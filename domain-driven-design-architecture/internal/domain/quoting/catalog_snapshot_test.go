package quoting

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/catalog"
)

func TestCatalogChangesDoNotMutateQuoteLineSnapshot(t *testing.T) {
	price, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Desk", catalog.ProductCategoryStandard, price, 30)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	productPrice, err := NewMoney(product.BasePrice().Cents(), product.BasePrice().Currency())
	if err != nil {
		t.Fatal(err)
	}
	line, err := NewQuoteLine(ProductSKU(product.SKU()), 2, productPrice)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	updatedPrice, err := catalog.NewPrice(17000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	if err := product.ChangeBasePrice(updatedPrice); err != nil {
		t.Fatal(err)
	}
	total, err := quote.Total()
	if err != nil {
		t.Fatal(err)
	}
	if total.Cents() != 30000 {
		t.Fatalf("quote total = %d, want 30000", total.Cents())
	}
}
