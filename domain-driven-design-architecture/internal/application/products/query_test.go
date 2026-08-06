package products

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/catalog"
)

func TestReaderProjectsProductAggregate(t *testing.T) {
	price, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Widget", catalog.ProductCategoryStandard, price, 30)
	if err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	reader.Save(product)
	details, err := reader.GetProduct("sku-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.Name != "Widget" || details.BasePriceCents != 15000 || !details.Active {
		t.Fatalf("unexpected details %+v", details)
	}
	active := true
	if got := reader.ListProducts(&active); len(got) != 1 || got[0].SKU != "sku-001" {
		t.Fatalf("unexpected products %+v", got)
	}
}
