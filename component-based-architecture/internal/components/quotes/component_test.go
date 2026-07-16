package quotes

import (
	"testing"

	"component-based-architecture/internal/components/customers"
	"component-based-architecture/internal/components/products"
)

func newQuoteComponent(t *testing.T) *Component {
	t.Helper()
	customerComponent := customers.NewComponent()
	if err := customerComponent.Register(customers.Customer{ID: "customer-001", Active: true}); err != nil {
		t.Fatalf("register customer: %v", err)
	}
	productComponent := products.NewComponent()
	if err := productComponent.Register(products.Product{SKU: "sku-001", Name: "Desk", Active: true, UnitPrice: 15000}); err != nil {
		t.Fatalf("register product: %v", err)
	}
	return NewComponent(customerComponent, productComponent)
}

func TestCreateDraftQuoteCreatesDraftForActiveCustomer(t *testing.T) {
	quoteComponent := newQuoteComponent(t)

	result, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", result.QuoteID)
	}
	if result.Status != QuoteStatusDraft {
		t.Fatalf("expected %s, got %s", QuoteStatusDraft, result.Status)
	}
}

func TestCreateDraftQuoteRejectsInactiveCustomer(t *testing.T) {
	customerComponent := customers.NewComponent()
	if err := customerComponent.Register(customers.Customer{ID: "customer-001", Active: false}); err != nil {
		t.Fatalf("register customer: %v", err)
	}
	quoteComponent := NewComponent(customerComponent, products.NewComponent())

	_, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != customers.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", customers.ErrCustomerInactive, err)
	}
}

func TestGetQuoteReturnsPublicDetails(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}

	result, err := quoteComponent.GetQuote(GetQuoteQuery{QuoteID: created.QuoteID})
	if err != nil {
		t.Fatalf("get quote: %v", err)
	}

	if result.QuoteID != created.QuoteID {
		t.Fatalf("expected quote id %s, got %s", created.QuoteID, result.QuoteID)
	}
	if result.CustomerID != "customer-001" {
		t.Fatalf("expected customer-001, got %s", result.CustomerID)
	}
	if result.Status != QuoteStatusDraft {
		t.Fatalf("expected %s, got %s", QuoteStatusDraft, result.Status)
	}
}

func TestGetQuoteReturnsNotFoundForUnknownQuote(t *testing.T) {
	quoteComponent := NewComponent(customers.NewComponent(), products.NewComponent())

	_, err := quoteComponent.GetQuote(GetQuoteQuery{QuoteID: "quote-999"})
	if err != ErrQuoteNotFound {
		t.Fatalf("expected %v, got %v", ErrQuoteNotFound, err)
	}
}

func TestAddQuoteLineAddsActiveProductToDraftQuote(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}

	result, err := quoteComponent.AddQuoteLine(AddQuoteLineCommand{
		QuoteID: created.QuoteID, ProductSKU: "sku-001", Quantity: 2,
	})
	if err != nil {
		t.Fatalf("add quote line: %v", err)
	}
	if result.LineCount != 1 {
		t.Fatalf("expected one line, got %d", result.LineCount)
	}

	details, err := quoteComponent.GetQuote(GetQuoteQuery{QuoteID: created.QuoteID})
	if err != nil {
		t.Fatalf("get quote: %v", err)
	}
	if details.LineCount != 1 {
		t.Fatalf("expected one line in details, got %d", details.LineCount)
	}
}
