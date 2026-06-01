package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubProductReader struct {
	product entities.Product
	err     error
}

func (g stubProductReader) FindBySKU(sku string) (entities.Product, error) {
	if g.err != nil {
		return entities.Product{}, g.err
	}

	return g.product, nil
}

type stubProductLister struct {
	products      []entities.Product
	err           error
	category      string
	availableOnly bool
}

func (g *stubProductLister) List(category string, availableOnly bool) ([]entities.Product, error) {
	g.category = category
	g.availableOnly = availableOnly
	if g.err != nil {
		return nil, g.err
	}

	return g.products, nil
}

type stubGetProductOutput struct {
	output GetProductOutput
}

func (o *stubGetProductOutput) Present(output GetProductOutput) error {
	o.output = output
	return nil
}

type stubListProductsOutput struct {
	output ListProductsOutput
}

func (o *stubListProductsOutput) Present(output ListProductsOutput) error {
	o.output = output
	return nil
}

func TestGetProductInteractorLoadsProduct(t *testing.T) {
	output := &stubGetProductOutput{}
	interactor := NewGetProductInteractor(stubProductReader{
		product: entities.Product{
			SKU:              "CHAIR-001",
			Name:             "Office Chair",
			Category:         "Standard",
			BasePrice:        10000,
			Available:        true,
			ReturnWindowDays: 30,
		},
	}, output)

	err := interactor.Execute(GetProductInput{SKU: "CHAIR-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.SKU != "CHAIR-001" {
		t.Fatalf("expected SKU CHAIR-001, got %s", output.output.SKU)
	}
}

func TestListProductsInteractorFiltersByCategoryAndAvailability(t *testing.T) {
	products := &stubProductLister{
		products: []entities.Product{
			{
				SKU:       "CHAIR-001",
				Name:      "Office Chair",
				Category:  "Standard",
				BasePrice: 10000,
				Available: true,
			},
			{
				SKU:       "DESK-001",
				Name:      "Standing Desk",
				Category:  "Standard",
				BasePrice: 25000,
				Available: true,
			},
		},
	}
	output := &stubListProductsOutput{}
	interactor := NewListProductsInteractor(products, output)

	err := interactor.Execute(ListProductsInput{Category: "Standard", AvailableOnly: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if products.category != "Standard" {
		t.Fatalf("expected category filter Standard, got %s", products.category)
	}

	if !products.availableOnly {
		t.Fatal("expected availableOnly filter to be true")
	}

	if output.output.Count != 2 {
		t.Fatalf("expected 2 products, got %d", output.output.Count)
	}
}
