package main

import (
	"math"
	"testing"
)

func TestConcurrentControllerFindsBestHypothesis(t *testing.T) {
	invoices := []Invoice{
		{ID: "INV-101", CustomerName: "Nikos Arvanitis", AmountCents: 12000},
		{ID: "INV-102", CustomerName: "Alexandros Papadopoulos", AmountCents: 45000},
		{ID: "INV-103", CustomerName: "Maria Georgiou", AmountCents: 45000},
	}
	payment := Payment{
		RawMemo:     "DEP BY A. PAPADOPOULOS REF 1042",
		AmountCents: 45000,
	}

	blackboard := NewBlackboard(payment, invoices)
	controller := NewController(0.7)
	controller.RegisterKS(ReferenceMatcher{})
	controller.RegisterKS(ExactAmountMatcher{})
	controller.RegisterKS(CustomerNameMatcher{})

	best, converged := controller.RunConcurrent(blackboard)

	if !converged {
		t.Fatal("expected the concurrent controller to converge")
	}
	if best.InvoiceID != "INV-102" {
		t.Fatalf("expected INV-102, got %s", best.InvoiceID)
	}
	if math.Abs(best.Score-0.7) > 1e-9 {
		t.Fatalf("expected score 0.7, got %.12f", best.Score)
	}
}
