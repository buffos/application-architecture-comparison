package application

import (
	"testing"

	"hexagonal-architecture/internal/adapters/repository/memory"
	"hexagonal-architecture/internal/core/domain"
)

func TestGetAndListProducts(t *testing.T) {
	productRepo := memory.NewProductRepository()
	getProduct := NewGetProductUseCase(productRepo)
	listProducts := NewListProductsUseCase(productRepo)

	_ = productRepo.Save(domain.Product{SKU: "CHAIR-001", Name: "Office Chair", Category: "Standard", BasePrice: 10000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "DESK-001", Name: "Executive Desk", Category: "CustomBuild", BasePrice: 50000, Available: true, ReturnWindowDays: 30})
	_ = productRepo.Save(domain.Product{SKU: "LAMP-001", Name: "Clearance Lamp", Category: "Clearance", BasePrice: 4000, Available: false, ReturnWindowDays: 0})

	product, err := getProduct.Execute("CHAIR-001")
	if err != nil {
		t.Fatalf("expected get product to succeed, got %v", err)
	}

	if product.Name != "Office Chair" {
		t.Fatalf("expected office chair, got %s", product.Name)
	}

	standardProducts, err := listProducts.Execute("Standard", false)
	if err != nil {
		t.Fatalf("expected list products to succeed, got %v", err)
	}

	if len(standardProducts) != 1 || standardProducts[0].SKU != "CHAIR-001" {
		t.Fatalf("expected one standard product CHAIR-001, got %+v", standardProducts)
	}

	availableProducts, err := listProducts.Execute("", true)
	if err != nil {
		t.Fatalf("expected list available products to succeed, got %v", err)
	}

	if len(availableProducts) != 2 {
		t.Fatalf("expected two available products, got %+v", availableProducts)
	}
}
