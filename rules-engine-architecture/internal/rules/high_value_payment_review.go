package rules

import (
	"fmt"

	"rules-engine-architecture/internal/engine"
)

// HighValuePaymentReviewRule requires review when the quote subtotal exceeds
// the configured threshold.
type HighValuePaymentReviewRule struct {
	thresholdCents int64
}

func NewHighValuePaymentReviewRule(thresholdCents int64) HighValuePaymentReviewRule {
	return HighValuePaymentReviewRule{thresholdCents: thresholdCents}
}

func (HighValuePaymentReviewRule) Name() string {
	return "High Value Payment Review Rule"
}

func (HighValuePaymentReviewRule) Priority() int {
	return 50
}

func (HighValuePaymentReviewRule) ConflictGroup() string {
	return ""
}

func (rule HighValuePaymentReviewRule) Evaluate(memory *engine.WorkingMemory) bool {
	return quoteSubtotal(memory) > rule.thresholdCents
}

func (rule HighValuePaymentReviewRule) Execute(memory *engine.WorkingMemory) error {
	subtotal := quoteSubtotal(memory)
	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: "payment-review",
		Message: fmt.Sprintf(
			"Quote subtotal of %s exceeds payment review threshold of %s",
			formatCents(subtotal),
			formatCents(rule.thresholdCents),
		),
	})
	return nil
}

func quoteSubtotal(memory *engine.WorkingMemory) int64 {
	var subtotal int64
	for _, line := range memory.Quote.Lines {
		subtotal += line.UnitPriceCents * int64(line.Quantity)
	}
	return subtotal
}

func formatCents(amountCents int64) string {
	return fmt.Sprintf("%d.%02d", amountCents/100, amountCents%100)
}
