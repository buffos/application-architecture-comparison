package products

import "testing"

type stubProductRepository struct {
	product Product
	err     error
}

func (r stubProductRepository) FindBySKU(sku string) (Product, error) {
	if r.err != nil {
		return Product{}, r.err
	}

	return r.product, nil
}

func (r stubProductRepository) Save(product Product) error {
	return nil
}

func TestGetProductForQuoteReturnsSellableProduct(t *testing.T) {
	service := NewService(stubProductRepository{
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
		t.Fatalf("expected no error, got %v", err)
	}

	if product.SKU != "sku-001" {
		t.Fatalf("expected sku-001, got %s", product.SKU)
	}
}

func TestGetProductForQuoteRejectsInactiveProduct(t *testing.T) {
	service := NewService(stubProductRepository{
		product: Product{
			SKU:    "sku-001",
			Active: false,
		},
	})

	_, err := service.GetProductForQuote("sku-001")
	if err != ErrProductInactive {
		t.Fatalf("expected %v, got %v", ErrProductInactive, err)
	}
}
