package scripts

import (
	"reflect"
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestCreateDraftQuotePersistsDraft(t *testing.T) {
	store := data.NewStore()
	store.Customers["customer-001"] = data.Customer{
		ID:     "customer-001",
		Active: true,
	}

	got, err := CreateDraftQuote(store, "customer-001")
	if err != nil {
		t.Fatalf("CreateDraftQuote() error = %v", err)
	}

	if got.ID != "quote-001" {
		t.Fatalf("quote ID = %q, want %q", got.ID, "quote-001")
	}
	if got.CustomerID != "customer-001" {
		t.Fatalf("customer ID = %q, want %q", got.CustomerID, "customer-001")
	}
	if got.Status != data.QuoteStatusDraft {
		t.Fatalf("status = %q, want %q", got.Status, data.QuoteStatusDraft)
	}

	saved, ok := store.Quotes[got.ID]
	if !ok {
		t.Fatalf("quote %q was not saved", got.ID)
	}
	if !reflect.DeepEqual(saved, got) {
		t.Fatalf("saved quote = %#v, want %#v", saved, got)
	}
}

func TestCreateDraftQuoteRejectsInactiveCustomer(t *testing.T) {
	store := data.NewStore()
	store.Customers["customer-001"] = data.Customer{
		ID:     "customer-001",
		Active: false,
	}

	_, err := CreateDraftQuote(store, "customer-001")
	if err != ErrCustomerInactive {
		t.Fatalf("error = %v, want %v", err, ErrCustomerInactive)
	}
	if len(store.Quotes) != 0 {
		t.Fatalf("quote count = %d, want 0", len(store.Quotes))
	}
}

func TestCreateDraftQuoteRejectsUnknownCustomer(t *testing.T) {
	store := data.NewStore()

	_, err := CreateDraftQuote(store, "customer-404")
	if err != ErrCustomerNotFound {
		t.Fatalf("error = %v, want %v", err, ErrCustomerNotFound)
	}
	if len(store.Quotes) != 0 {
		t.Fatalf("quote count = %d, want 0", len(store.Quotes))
	}
}
