package quotes

import (
	"testing"

	"modular-monolith/internal/modules/approvals"
	"modular-monolith/internal/modules/customers"
	"modular-monolith/internal/modules/products"
)

type stubQuoteRepository struct {
	saved Quote
}

func (r *stubQuoteRepository) Save(quote Quote) error {
	r.saved = quote
	return nil
}

func (r *stubQuoteRepository) FindByID(id string) (Quote, error) {
	return r.saved, nil
}

func (r *stubQuoteRepository) ListByStatus(status string) ([]Quote, error) {
	if r.saved.ID == "" {
		return nil, nil
	}
	if status == "" || r.saved.Status == status {
		return []Quote{r.saved}, nil
	}
	return nil, nil
}

type stubCustomerDirectory struct {
	err error
}

func (d stubCustomerDirectory) RequireActiveCustomer(id string) error {
	return d.err
}

type stubProductCatalog struct {
	product products.ProductForQuote
	err     error
}

func (c stubProductCatalog) GetProductForQuote(sku string) (products.ProductForQuote, error) {
	if c.err != nil {
		return products.ProductForQuote{}, c.err
	}

	return c.product, nil
}

type stubApprovalEvaluator struct {
	requiresApproval bool
}

type stubQuotePricer struct {
	unitPrice int
	err       error
}

func (e stubApprovalEvaluator) RequiresApproval(submission approvals.QuoteSubmission) bool {
	return e.requiresApproval
}

func (p stubQuotePricer) UnitPrice(product products.ProductForQuote) (int, error) {
	if p.err != nil {
		return 0, p.err
	}
	if p.unitPrice == 0 {
		return product.UnitPrice, nil
	}
	return p.unitPrice, nil
}

func TestCreateDraftQuoteCreatesDraftForActiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{})

	result, err := service.CreateDraftQuote(CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != QuoteStatusDraft {
		t.Fatalf("expected status %s, got %s", QuoteStatusDraft, result.Status)
	}

	if quotes.saved.CustomerID != "customer-001" {
		t.Fatalf("expected customer-001, got %s", quotes.saved.CustomerID)
	}
}

func TestCreateDraftQuoteRejectsInactiveCustomer(t *testing.T) {
	quotes := &stubQuoteRepository{}
	service := NewService(quotes, stubCustomerDirectory{
		err: customers.ErrCustomerInactive,
	}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{})

	_, err := service.CreateDraftQuote(CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != customers.ErrCustomerInactive {
		t.Fatalf("expected %v, got %v", customers.ErrCustomerInactive, err)
	}
}

func TestAddQuoteLineAddsLineToExistingQuote(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{
		product: products.ProductForQuote{
			SKU:       "sku-001",
			Name:      "Desk",
			Category:  "Standard",
			UnitPrice: 15000,
		},
	}, stubQuotePricer{}, stubApprovalEvaluator{})

	result, err := service.AddQuoteLine(AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.LineCount != 1 {
		t.Fatalf("expected line count 1, got %d", result.LineCount)
	}

	if result.TotalItems != 2 {
		t.Fatalf("expected total items 2, got %d", result.TotalItems)
	}
}

func TestAddQuoteLineUsesPricingPolicyUnitPrice(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{
		product: products.ProductForQuote{
			SKU:       "sku-001",
			Name:      "Desk",
			Category:  "Standard",
			UnitPrice: 10000,
		},
	}, stubQuotePricer{unitPrice: 9500}, stubApprovalEvaluator{})

	_, err := service.AddQuoteLine(AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quotes.saved.Lines[0].UnitPrice != 9500 {
		t.Fatalf("expected priced unit 9500, got %d", quotes.saved.Lines[0].UnitPrice)
	}
}

func TestSubmitQuoteSubmitsQuoteWithLines(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
			Lines: []QuoteLine{
				{ProductSKU: "sku-001", Quantity: 1, UnitPrice: 15000},
			},
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{
		requiresApproval: false,
	})

	result, err := service.SubmitQuote(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != QuoteStatusApproved {
		t.Fatalf("expected status %s, got %s", QuoteStatusApproved, result.Status)
	}
}

func TestSubmitQuoteRejectsEmptyQuote(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{})

	_, err := service.SubmitQuote(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != ErrQuoteCannotBeSubmittedWithoutLines {
		t.Fatalf("expected %v, got %v", ErrQuoteCannotBeSubmittedWithoutLines, err)
	}
}

func TestSubmitQuoteMovesCustomQuoteToPendingApproval(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
			Lines: []QuoteLine{
				{ProductSKU: "sku-002", ProductCategory: "CustomBuild", Quantity: 1, UnitPrice: 45000},
			},
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{
		requiresApproval: true,
	})

	result, err := service.SubmitQuote(SubmitQuoteCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != QuoteStatusPendingApproval {
		t.Fatalf("expected status %s, got %s", QuoteStatusPendingApproval, result.Status)
	}
}

func TestApproveQuoteApprovesPendingQuote(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusPendingApproval,
			Lines: []QuoteLine{
				{ProductSKU: "sku-002", ProductCategory: "CustomBuild", Quantity: 1, UnitPrice: 45000},
			},
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{})

	result, err := service.ApproveQuote(ApproveQuoteCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != QuoteStatusApproved {
		t.Fatalf("expected status %s, got %s", QuoteStatusApproved, result.Status)
	}
}

func TestApproveQuoteRejectsAlreadyApprovedQuote(t *testing.T) {
	quotes := &stubQuoteRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusApproved,
			Lines: []QuoteLine{
				{ProductSKU: "sku-001", ProductCategory: "Standard", Quantity: 1, UnitPrice: 15000},
			},
		},
	}
	service := NewService(quotes, stubCustomerDirectory{}, stubProductCatalog{}, stubQuotePricer{}, stubApprovalEvaluator{})

	_, err := service.ApproveQuote(ApproveQuoteCommand{QuoteID: "quote-001"})
	if err != ErrQuoteNotApprovable {
		t.Fatalf("expected %v, got %v", ErrQuoteNotApprovable, err)
	}
}
