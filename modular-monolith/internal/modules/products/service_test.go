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

func (r stubProductRepository) List(category string, activeOnly bool) ([]Product, error) {
	list := []Product{}
	if r.product.SKU == "" {
		return list, nil
	}
	if category != "" && r.product.Category != category {
		return list, nil
	}
	if activeOnly && !r.product.Active {
		return list, nil
	}
	return []Product{r.product}, nil
}

func (r stubProductRepository) Save(product Product) error {
	return nil
}

func TestGetProductForQuoteReturnsSellableProduct(t *testing.T) {
	service := NewService(stubProductRepository{
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
		t.Fatalf("expected no error, got %v", err)
	}

	if product.SKU != "sku-001" {
		t.Fatalf("expected sku-001, got %s", product.SKU)
	}

	if product.ReturnWindowDays != 30 {
		t.Fatalf("expected return window 30, got %d", product.ReturnWindowDays)
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

func TestGetProductReturnsStoredProduct(t *testing.T) {
	service := NewService(stubProductRepository{
		product: Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			Active:           true,
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	})

	product, err := service.GetProduct(GetProductQuery{SKU: "sku-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.SKU != "sku-001" || !product.Active {
		t.Fatalf("expected stored product details, got %+v", product)
	}
}

func TestListProductsFiltersByCategoryAndActivity(t *testing.T) {
	service := NewService(stubListProductRepository{
		products: map[string]Product{
			"sku-001": {
				SKU:              "sku-001",
				Name:             "Desk",
				Category:         "Standard",
				Active:           true,
				UnitPrice:        15000,
				ReturnWindowDays: 30,
			},
			"sku-002": {
				SKU:              "sku-002",
				Name:             "Custom Desk",
				Category:         "CustomBuild",
				Active:           true,
				UnitPrice:        45000,
				ReturnWindowDays: 14,
			},
			"sku-003": {
				SKU:              "sku-003",
				Name:             "Old Chair",
				Category:         "Standard",
				Active:           false,
				UnitPrice:        5000,
				ReturnWindowDays: 30,
			},
		},
	})

	result, err := service.ListProducts(ListProductsQuery{Category: "Standard", ActiveOnly: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].SKU != "sku-001" {
		t.Fatalf("expected one active standard product, got %+v", result)
	}
}

type stubListProductRepository struct {
	products map[string]Product
}

func (r stubListProductRepository) FindBySKU(sku string) (Product, error) {
	product, ok := r.products[sku]
	if !ok {
		return Product{}, ErrProductNotFound
	}
	return product, nil
}

func (r stubListProductRepository) List(category string, activeOnly bool) ([]Product, error) {
	list := make([]Product, 0, len(r.products))
	for _, product := range r.products {
		if category != "" && product.Category != category {
			continue
		}
		if activeOnly && !product.Active {
			continue
		}
		list = append(list, product)
	}
	return list, nil
}

func (r stubListProductRepository) Save(product Product) error {
	return nil
}
