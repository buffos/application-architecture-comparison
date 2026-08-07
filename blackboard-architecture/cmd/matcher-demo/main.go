package main

import (
	"fmt"
	"strings"
)

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

func (bb *Blackboard) HasConverged(threshold float64) bool {
	best, found := bb.BestHypothesis()
	return found && best.Score >= threshold
}

type KnowledgeSource interface {
	Name() string
	Execute(bb *Blackboard)
}

type ExactAmountMatcher struct{}

func (ExactAmountMatcher) Name() string {
	return "Exact Amount Matcher"
}

func (ExactAmountMatcher) Execute(bb *Blackboard) {
	for _, invoice := range bb.Invoices {
		if invoice.AmountCents == bb.Payment.AmountCents {
			bb.AddEvidence(
				invoice.ID,
				0.4,
				fmt.Sprintf("Exact amount match for %s (+0.4)", formatCents(invoice.AmountCents)),
			)
		}
	}
}

type ReferenceMatcher struct{}

func (ReferenceMatcher) Name() string {
	return "Invoice Reference Matcher"
}

func (ReferenceMatcher) Execute(bb *Blackboard) {
	memo := strings.ToLower(bb.Payment.RawMemo)

	for _, invoice := range bb.Invoices {
		if strings.Contains(memo, strings.ToLower(invoice.ID)) {
			bb.AddEvidence(
				invoice.ID,
				0.5,
				fmt.Sprintf("Found invoice reference %q in memo (+0.5)", invoice.ID),
			)
		}
	}
}

type CustomerNameMatcher struct{}

func (CustomerNameMatcher) Name() string {
	return "Customer Name Matcher"
}

func (CustomerNameMatcher) Execute(bb *Blackboard) {
	memo := strings.ToLower(bb.Payment.RawMemo)

	for _, invoice := range bb.Invoices {
		parts := strings.Fields(invoice.CustomerName)
		if len(parts) == 0 {
			continue
		}

		lastName := parts[len(parts)-1]
		if strings.Contains(memo, strings.ToLower(lastName)) {
			bb.AddEvidence(
				invoice.ID,
				0.3,
				fmt.Sprintf("Found customer surname %q in memo (+0.3)", lastName),
			)
		}
	}
}

type Controller struct {
	sources             []KnowledgeSource
	confidenceThreshold float64
}

func NewController(confidenceThreshold float64) *Controller {
	return &Controller{
		sources:             []KnowledgeSource{},
		confidenceThreshold: confidenceThreshold,
	}
}

func (c *Controller) RegisterKS(source KnowledgeSource) {
	c.sources = append(c.sources, source)
}

func (c *Controller) Run(bb *Blackboard) (MatchHypothesis, bool) {
	for index, source := range c.sources {
		if bb.HasConverged(c.confidenceThreshold) {
			best, _ := bb.BestHypothesis()
			fmt.Printf("Converged at %.1f; skipped %d source(s)\n", best.Score, len(c.sources)-index)
			break
		}

		fmt.Printf("Controller executing: %s\n", source.Name())
		source.Execute(bb)
	}

	best, found := bb.BestHypothesis()
	return best, found && best.Score >= c.confidenceThreshold
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
	controller := NewController(0.7)
	controller.RegisterKS(ReferenceMatcher{})
	controller.RegisterKS(ExactAmountMatcher{})
	controller.RegisterKS(CustomerNameMatcher{})
	best, converged := controller.Run(blackboard)

	fmt.Printf("Payment memo: %q\n", blackboard.Payment.RawMemo)
	fmt.Printf("Payment amount: %s\n", formatCents(blackboard.Payment.AmountCents))
	fmt.Println("Candidate invoices:")

	for _, invoice := range blackboard.Invoices {
		hypothesis := blackboard.Hypotheses[invoice.ID]
		fmt.Printf("- %s, %s, amount %s, score %.1f\n",
			invoice.ID,
			invoice.CustomerName,
			formatCents(invoice.AmountCents),
			hypothesis.Score,
		)
		for _, reason := range hypothesis.Reasons {
			fmt.Printf("  reason: %s\n", reason)
		}
	}

	fmt.Printf("Controller result: %s with score %.1f\n", best.InvoiceID, best.Score)
	fmt.Printf("Converged: %t\n", converged)
}
