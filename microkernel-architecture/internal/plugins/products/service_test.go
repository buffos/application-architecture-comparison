package products

import "testing"

type stubRepository struct {
	product Product
	err     error
}

func (r stubRepository) FindBySKU(sku string) (Product, error) {
	return r.product, r.err
}

func (r stubRepository) Save(product Product) error {
	return nil
}

func TestGetProductForQuote(t *testing.T) {
	service := NewService(stubRepository{
		product: Product{
			SKU:       "sku-001",
			Name:      "Desk",
			Category:  "Standard",
			Active:    true,
			UnitPrice: 15000,
		},
	})

	product, err := service.GetProductForQuote("sku-001")
	if err != nil {
		t.Fatalf("expected product lookup to succeed, got %v", err)
	}

	if product.SKU != "sku-001" {
		t.Fatalf("expected sku sku-001, got %s", product.SKU)
	}
}
