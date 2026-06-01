package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubProductFinder struct {
	product domain.Product
	list    []domain.Product
	err     error
}

func (l stubProductFinder) FindBySKU(sku string) (domain.Product, error) {
	if l.err != nil {
		return domain.Product{}, l.err
	}

	return l.product, nil
}

func (l stubProductFinder) List(category string, activeOnly bool) ([]domain.Product, error) {
	if l.err != nil {
		return nil, l.err
	}

	result := make([]domain.Product, 0)
	for _, product := range l.list {
		if category != "" && product.Category != category {
			continue
		}

		if activeOnly && !product.Active {
			continue
		}

		result = append(result, product)
	}

	return result, nil
}

func TestGetProductServiceReturnsDetails(t *testing.T) {
	service := NewGetProductService(stubProductFinder{
		product: domain.Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			Active:           true,
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	})

	result, err := service.Execute(GetProductQuery{SKU: "sku-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.SKU != "sku-001" {
		t.Fatalf("expected sku-001, got %s", result.SKU)
	}
}

func TestListProductsServiceFiltersByCategoryAndAvailability(t *testing.T) {
	service := NewListProductsService(stubProductFinder{
		list: []domain.Product{
			{SKU: "sku-001", Category: "Standard", Active: true},
			{SKU: "sku-002", Category: "CustomBuild", Active: true},
			{SKU: "sku-003", Category: "Standard", Active: false},
		},
	})

	result, err := service.Execute(ListProductsQuery{
		Category:   "Standard",
		ActiveOnly: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].SKU != "sku-001" {
		t.Fatalf("expected sku-001, got %s", result[0].SKU)
	}
}
