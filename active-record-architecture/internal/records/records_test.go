package records

import "testing"

func TestCustomerSaveAndFindRoundTrip(t *testing.T) {
	db := NewDatabase()
	want := NewCustomer(db, "customer-001", true)

	if err := want.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}

	got, err := FindCustomer(db, want.ID)
	if err != nil {
		t.Fatalf("FindCustomer() error = %v", err)
	}
	if got.ID != want.ID || got.Active != want.Active {
		t.Fatalf("loaded customer = %#v, want id=%q active=%t", got, want.ID, want.Active)
	}
}

func TestNewDraftQuoteRequiresExplicitActiveRecordSave(t *testing.T) {
	db := NewDatabase()
	customer := NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}

	quote, err := NewDraftQuote(db, customer.ID)
	if err != nil {
		t.Fatalf("NewDraftQuote() error = %v", err)
	}
	if quote.ID != "quote-001" {
		t.Fatalf("quote ID = %q, want %q", quote.ID, "quote-001")
	}
	if quote.Status != QuoteStatusDraft {
		t.Fatalf("quote status = %q, want %q", quote.Status, QuoteStatusDraft)
	}

	if _, err := FindQuote(db, quote.ID); err != ErrQuoteNotFound {
		t.Fatalf("FindQuote() before Save error = %v, want %v", err, ErrQuoteNotFound)
	}

	if err := quote.Save(); err != nil {
		t.Fatalf("Quote.Save() error = %v", err)
	}

	saved, err := FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() after Save error = %v", err)
	}
	if saved.CustomerID != customer.ID || saved.Status != QuoteStatusDraft {
		t.Fatalf("saved quote = %#v, want customer=%q status=%q", saved, customer.ID, QuoteStatusDraft)
	}
}

func TestNewDraftQuoteRejectsInactiveCustomer(t *testing.T) {
	db := NewDatabase()
	customer := NewCustomer(db, "customer-001", false)
	if err := customer.Save(); err != nil {
		t.Fatalf("Customer.Save() error = %v", err)
	}

	_, err := NewDraftQuote(db, customer.ID)
	if err != ErrCustomerInactive {
		t.Fatalf("NewDraftQuote() error = %v, want %v", err, ErrCustomerInactive)
	}
}

func TestNewDraftQuoteRejectsUnknownCustomer(t *testing.T) {
	db := NewDatabase()

	_, err := NewDraftQuote(db, "customer-404")
	if err != ErrCustomerNotFound {
		t.Fatalf("NewDraftQuote() error = %v, want %v", err, ErrCustomerNotFound)
	}
}
