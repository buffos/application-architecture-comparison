package products

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

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

func (r stubRepository) List(category string, active *bool) ([]Product, error) {
	if r.err != nil {
		return nil, r.err
	}

	if category != "" && r.product.Category != category {
		return []Product{}, nil
	}

	if active != nil && r.product.Active != *active {
		return []Product{}, nil
	}

	return []Product{r.product}, nil
}

func TestGetProductForQuote(t *testing.T) {
	service := NewService(stubRepository{
		product: Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			Active:           true,
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	})

	product, err := service.GetProductForQuote("sku-001")
	if err != nil {
		t.Fatalf("expected product lookup to succeed, got %v", err)
	}

	if product.SKU != "sku-001" {
		t.Fatalf("expected sku sku-001, got %s", product.SKU)
	}

	if product.ReturnWindowDays != 30 {
		t.Fatalf("expected return window 30, got %d", product.ReturnWindowDays)
	}
}

func TestGetProduct(t *testing.T) {
	service := NewService(stubRepository{
		product: Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			Active:           true,
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	})

	product, err := service.GetProduct(kernel.GetProductQuery{SKU: "sku-001"})
	if err != nil {
		t.Fatalf("expected product query to succeed, got %v", err)
	}

	if product.SKU != "sku-001" || !product.Active {
		t.Fatalf("unexpected product details %+v", product)
	}
}

func TestListProducts(t *testing.T) {
	active := true
	service := NewService(stubRepository{
		product: Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			Active:           true,
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	})

	products, err := service.ListProducts(kernel.ListProductsQuery{
		Category: "Standard",
		Active:   &active,
	})
	if err != nil {
		t.Fatalf("expected product list to succeed, got %v", err)
	}

	if len(products) != 1 || products[0].SKU != "sku-001" {
		t.Fatalf("unexpected product list %+v", products)
	}
}
