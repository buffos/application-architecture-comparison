package main

import "fmt"

// Invoice is one unpaid invoice that may match the incoming payment.
type Invoice struct {
	ID           string
	CustomerName string
	AmountCents  int64
}

// Payment is the unstructured bank input we want to explain.
type Payment struct {
	RawMemo     string
	AmountCents int64
}

// MatchHypothesis is a candidate explanation for the payment.
// The score and reasons will be populated by Knowledge Sources in later lessons.
type MatchHypothesis struct {
	InvoiceID string
	Score     float64
	Reasons   []string
}

// Blackboard is the shared working memory of the matcher.
//
// For this first lesson it is deliberately a plain in-memory structure.
// Concurrency protection will be introduced only when parallel Knowledge
// Sources are added.
type Blackboard struct {
	Payment    Payment
	Invoices   []Invoice
	Hypotheses map[string]MatchHypothesis
}

func NewBlackboard(payment Payment, invoices []Invoice) *Blackboard {
	hypotheses := make(map[string]MatchHypothesis, len(invoices))
	for _, invoice := range invoices {
		hypotheses[invoice.ID] = MatchHypothesis{
			InvoiceID: invoice.ID,
			Reasons:   []string{},
		}
	}

	return &Blackboard{
		Payment:    payment,
		Invoices:   invoices,
		Hypotheses: hypotheses,
	}
}

func (bb *Blackboard) AddEvidence(invoiceID string, points float64, reason string) {
	hypothesis, exists := bb.Hypotheses[invoiceID]
	if !exists {
		return
	}

	hypothesis.Score += points
	if hypothesis.Score > 1.0 {
		hypothesis.Score = 1.0
	}
	hypothesis.Reasons = append(hypothesis.Reasons, reason)
	bb.Hypotheses[invoiceID] = hypothesis
}

func (bb *Blackboard) BestHypothesis() (MatchHypothesis, bool) {
	var best MatchHypothesis
	found := false

	for _, hypothesis := range bb.Hypotheses {
		if !found || hypothesis.Score > best.Score {
			best = hypothesis
			found = true
		}
	}

	return best, found
}

func formatCents(amountCents int64) string {
	return fmt.Sprintf("%d.%02d", amountCents/100, amountCents%100)
}

func main() {
	unpaidInvoices := []Invoice{
		{ID: "INV-101", CustomerName: "Nikos Arvanitis", AmountCents: 12000},
		{ID: "INV-102", CustomerName: "Alexandros Papadopoulos", AmountCents: 45000},
		{ID: "INV-103", CustomerName: "Maria Georgiou", AmountCents: 45000},
	}

	incomingPayment := Payment{
		RawMemo:     "DEP BY A. PAPADOPOULOS REF 1042",
		AmountCents: 45000,
	}

	blackboard := NewBlackboard(incomingPayment, unpaidInvoices)
	blackboard.AddEvidence(
		"INV-102",
		0.4,
		"Temporary exact amount observation (+0.4)",
	)

	fmt.Printf("Payment memo: %q\n", blackboard.Payment.RawMemo)
	fmt.Printf("Payment amount: %s\n", formatCents(blackboard.Payment.AmountCents))
	fmt.Println("Candidate invoices:")

	for _, invoice := range blackboard.Invoices {
		hypothesis := blackboard.Hypotheses[invoice.ID]
		fmt.Printf("- %s, %s, amount %s, initial score %.1f\n",
			invoice.ID,
			invoice.CustomerName,
			formatCents(invoice.AmountCents),
			hypothesis.Score,
		)
	}

	if best, ok := blackboard.BestHypothesis(); ok {
		fmt.Printf("Best hypothesis: %s with score %.1f\n", best.InvoiceID, best.Score)
		fmt.Println("Reasons:")
		for _, reason := range best.Reasons {
			fmt.Printf("- %s\n", reason)
		}
	}
}
