package quotes

import (
	"testing"

	"component-based-architecture/internal/components/approvals"
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
	if err := productComponent.Register(products.Product{SKU: "sku-001", Name: "Desk", Category: "Standard", Active: true, UnitPrice: 15000}); err != nil {
		t.Fatalf("register product: %v", err)
	}
	if err := productComponent.Register(products.Product{SKU: "sku-002", Name: "Custom Desk", Category: "CustomBuild", Active: true, UnitPrice: 45000}); err != nil {
		t.Fatalf("register custom product: %v", err)
	}
	return NewComponent(customerComponent, productComponent, approvals.NewComponent())
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
	quoteComponent := NewComponent(customerComponent, products.NewComponent(), approvals.NewComponent())

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
	quoteComponent := NewComponent(customers.NewComponent(), products.NewComponent(), approvals.NewComponent())

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

func TestSubmitQuoteApprovesDraftWithLines(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}
	if _, err := quoteComponent.AddQuoteLine(AddQuoteLineCommand{QuoteID: created.QuoteID, ProductSKU: "sku-001", Quantity: 1}); err != nil {
		t.Fatalf("add quote line: %v", err)
	}

	result, err := quoteComponent.SubmitQuote(SubmitQuoteCommand{QuoteID: created.QuoteID})
	if err != nil {
		t.Fatalf("submit quote: %v", err)
	}
	if result.Status != QuoteStatusApproved {
		t.Fatalf("expected %s, got %s", QuoteStatusApproved, result.Status)
	}
}

func TestSubmitQuoteRejectsEmptyDraft(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}

	_, err = quoteComponent.SubmitQuote(SubmitQuoteCommand{QuoteID: created.QuoteID})
	if err != ErrQuoteCannotBeSubmittedWithoutLines {
		t.Fatalf("expected %v, got %v", ErrQuoteCannotBeSubmittedWithoutLines, err)
	}
}

func TestAddQuoteLineRejectsApprovedQuote(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}
	if _, err := quoteComponent.AddQuoteLine(AddQuoteLineCommand{QuoteID: created.QuoteID, ProductSKU: "sku-001", Quantity: 1}); err != nil {
		t.Fatalf("add quote line: %v", err)
	}
	if _, err := quoteComponent.SubmitQuote(SubmitQuoteCommand{QuoteID: created.QuoteID}); err != nil {
		t.Fatalf("submit quote: %v", err)
	}

	_, err = quoteComponent.AddQuoteLine(AddQuoteLineCommand{QuoteID: created.QuoteID, ProductSKU: "sku-001", Quantity: 1})
	if err != ErrQuoteNotEditable {
		t.Fatalf("expected %v, got %v", ErrQuoteNotEditable, err)
	}
}

func TestSubmitQuoteSendsCustomBuildToPendingApproval(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}
	if _, err := quoteComponent.AddQuoteLine(AddQuoteLineCommand{QuoteID: created.QuoteID, ProductSKU: "sku-002", Quantity: 1}); err != nil {
		t.Fatalf("add custom quote line: %v", err)
	}

	result, err := quoteComponent.SubmitQuote(SubmitQuoteCommand{QuoteID: created.QuoteID})
	if err != nil {
		t.Fatalf("submit quote: %v", err)
	}
	if result.Status != QuoteStatusPendingApproval {
		t.Fatalf("expected %s, got %s", QuoteStatusPendingApproval, result.Status)
	}
}

func TestApproveQuoteApprovesPendingQuote(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}
	if _, err := quoteComponent.AddQuoteLine(AddQuoteLineCommand{QuoteID: created.QuoteID, ProductSKU: "sku-002", Quantity: 1}); err != nil {
		t.Fatalf("add custom quote line: %v", err)
	}
	if _, err := quoteComponent.SubmitQuote(SubmitQuoteCommand{QuoteID: created.QuoteID}); err != nil {
		t.Fatalf("submit quote: %v", err)
	}

	result, err := quoteComponent.ApproveQuote(ApproveQuoteCommand{QuoteID: created.QuoteID})
	if err != nil {
		t.Fatalf("approve quote: %v", err)
	}
	if result.Status != QuoteStatusApproved {
		t.Fatalf("expected %s, got %s", QuoteStatusApproved, result.Status)
	}
}

func TestApproveQuoteRejectsAlreadyApprovedQuote(t *testing.T) {
	quoteComponent := newQuoteComponent(t)
	created, err := quoteComponent.CreateDraftQuote(CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		t.Fatalf("create draft quote: %v", err)
	}
	if _, err := quoteComponent.AddQuoteLine(AddQuoteLineCommand{QuoteID: created.QuoteID, ProductSKU: "sku-002", Quantity: 1}); err != nil {
		t.Fatalf("add custom quote line: %v", err)
	}
	if _, err := quoteComponent.SubmitQuote(SubmitQuoteCommand{QuoteID: created.QuoteID}); err != nil {
		t.Fatalf("submit quote: %v", err)
	}
	if _, err := quoteComponent.ApproveQuote(ApproveQuoteCommand{QuoteID: created.QuoteID}); err != nil {
		t.Fatalf("approve quote: %v", err)
	}

	_, err = quoteComponent.ApproveQuote(ApproveQuoteCommand{QuoteID: created.QuoteID})
	if err != ErrQuoteNotApprovable {
		t.Fatalf("expected %v, got %v", ErrQuoteNotApprovable, err)
	}
}
