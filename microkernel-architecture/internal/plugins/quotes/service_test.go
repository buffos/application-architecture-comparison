package quotes

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved Quote
}

func (r *stubRepository) FindByID(id string) (Quote, error) {
	if r.saved.ID == id {
		return r.saved, nil
	}

	return Quote{}, ErrQuoteNotFound
}

func (r *stubRepository) ListByStatus(status string) ([]Quote, error) {
	if r.saved.ID == "" {
		return []Quote{}, nil
	}

	if status == "" || r.saved.Status == status {
		return []Quote{r.saved}, nil
	}

	return []Quote{}, nil
}

func (r *stubRepository) Save(quote Quote) error {
	r.saved = quote
	return nil
}

type stubCustomerDirectory struct {
	err error
}

func (d stubCustomerDirectory) RequireActiveCustomer(id string) error {
	return d.err
}

type stubProductCatalog struct {
	product kernel.Product
	err     error
}

func (c stubProductCatalog) GetProductForQuote(sku string) (kernel.Product, error) {
	return c.product, c.err
}

type stubApprovalPolicy struct {
	requiresApproval bool
}

func (p stubApprovalPolicy) RequiresApproval(submission kernel.QuoteSubmission) bool {
	return p.requiresApproval
}

func TestCreateDraftQuote(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	result, err := service.CreateDraftQuote(kernel.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		t.Fatalf("expected create draft quote to succeed, got %v", err)
	}

	if result.CustomerID != "customer-001" {
		t.Fatalf("expected customer id to be preserved, got %s", result.CustomerID)
	}

	if repository.saved.Status != QuoteStatusDraft {
		t.Fatalf("expected saved quote status %s, got %s", QuoteStatusDraft, repository.saved.Status)
	}
}

func TestGetQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	result, err := service.GetQuote(kernel.GetQuoteQuery{
		QuoteID: "quote-001",
	})
	if err != nil {
		t.Fatalf("expected get quote to succeed, got %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote id quote-001, got %s", result.QuoteID)
	}

	if result.TotalAmount != 0 {
		t.Fatalf("expected empty quote total amount 0, got %d", result.TotalAmount)
	}
}

func TestListQuotesByStatus(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusApproved,
			Lines: []QuoteLine{
				{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000},
			},
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	result, err := service.ListQuotes(kernel.ListQuotesQuery{
		Status: QuoteStatusApproved,
	})
	if err != nil {
		t.Fatalf("expected list quotes to succeed, got %v", err)
	}

	if len(result) != 1 || result[0].QuoteID != "quote-001" {
		t.Fatalf("unexpected quote list %+v", result)
	}

	if result[0].TotalAmount != 30000 {
		t.Fatalf("expected total amount 30000, got %d", result[0].TotalAmount)
	}
}

func TestAddQuoteLine(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{
		product: kernel.Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	}, stubApprovalPolicy{})

	result, err := service.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("expected add quote line to succeed, got %v", err)
	}

	if result.LineCount != 1 {
		t.Fatalf("expected one line, got %d", result.LineCount)
	}

	if repository.saved.TotalQuantity() != 2 {
		t.Fatalf("expected total quantity 2, got %d", repository.saved.TotalQuantity())
	}
}

func TestSubmitQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
			Lines: []QuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	result, err := service.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: "quote-001",
	})
	if err != nil {
		t.Fatalf("expected submit quote to succeed, got %v", err)
	}

	if result.Status != QuoteStatusApproved {
		t.Fatalf("expected approved status, got %s", result.Status)
	}
}

func TestSubmitQuoteRejectsEmptyQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	_, err := service.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: "quote-001",
	})
	if err != ErrQuoteCannotBeSubmittedWithoutLines {
		t.Fatalf("expected empty quote submit error, got %v", err)
	}
}

func TestAddQuoteLineRejectsSubmittedQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusApproved,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{
		product: kernel.Product{
			SKU:              "sku-001",
			Name:             "Desk",
			Category:         "Standard",
			UnitPrice:        15000,
			ReturnWindowDays: 30,
		},
	}, stubApprovalPolicy{})

	_, err := service.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    "quote-001",
		ProductSKU: "sku-001",
		Quantity:   1,
	})
	if err != ErrQuoteNotEditable {
		t.Fatalf("expected not editable error, got %v", err)
	}
}

func TestSubmitQuoteCanRequireApproval(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusDraft,
			Lines: []QuoteLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{
		requiresApproval: true,
	})

	result, err := service.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: "quote-001",
	})
	if err != nil {
		t.Fatalf("expected submit quote to succeed, got %v", err)
	}

	if result.Status != QuoteStatusPendingApproval {
		t.Fatalf("expected pending approval status, got %s", result.Status)
	}
}

func TestApproveQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusPendingApproval,
			Lines: []QuoteLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	result, err := service.ApproveQuote(kernel.ApproveQuoteCommand{
		QuoteID: "quote-001",
	})
	if err != nil {
		t.Fatalf("expected approve quote to succeed, got %v", err)
	}

	if result.Status != QuoteStatusApproved {
		t.Fatalf("expected approved status, got %s", result.Status)
	}
}

func TestApproveQuoteRejectsNonPendingQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusApproved,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	_, err := service.ApproveQuote(kernel.ApproveQuoteCommand{
		QuoteID: "quote-001",
	})
	if err != ErrQuoteNotApprovable {
		t.Fatalf("expected not approvable error, got %v", err)
	}
}

func TestGetApprovedQuoteForOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusApproved,
			Lines: []QuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	result, err := service.GetApprovedQuoteForOrder("quote-001")
	if err != nil {
		t.Fatalf("expected approved quote lookup to succeed, got %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote id quote-001, got %s", result.QuoteID)
	}
}

func TestGetApprovedQuoteForOrderRejectsNonApprovedQuote(t *testing.T) {
	repository := &stubRepository{
		saved: Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     QuoteStatusPendingApproval,
		},
	}
	service := NewService(repository, stubCustomerDirectory{}, stubProductCatalog{}, stubApprovalPolicy{})

	_, err := service.GetApprovedQuoteForOrder("quote-001")
	if err != ErrQuoteNotConvertible {
		t.Fatalf("expected not convertible error, got %v", err)
	}
}
