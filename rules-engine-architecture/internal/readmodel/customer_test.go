package readmodel

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestProjectCustomersSortsAndCopiesCustomerFacts(t *testing.T) {
	customers := []engine.CustomerFact{
		{ID: "CUST-002", Name: "Maria Georgiou", Tier: "Standard"},
		{ID: "CUST-001", Name: "Alexandros Papadopoulos", Tier: "Preferred", InvoiceTerms: true},
	}

	views := ProjectCustomers(customers)

	if len(views) != 2 || views[0].ID != "CUST-001" || views[1].ID != "CUST-002" {
		t.Fatalf("expected sorted customer views, got %+v", views)
	}
	if views[0].Tier != "Preferred" || !views[0].InvoiceTerms {
		t.Fatalf("expected customer commercial facts, got %+v", views[0])
	}

	customers[0].Name = "Changed"
	if views[1].Name != "Maria Georgiou" {
		t.Fatal("expected customer view to copy source values")
	}
}
