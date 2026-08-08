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
	return quoteSubtotal(memory) > rule.thresholdCents &&
		memory.PaymentReview.Status != engine.PaymentReviewApproved
}

func (rule HighValuePaymentReviewRule) Execute(memory *engine.WorkingMemory) error {
	subtotal := quoteSubtotal(memory)
	severity := "payment-review"
	message := fmt.Sprintf(
		"Quote subtotal of %s exceeds payment review threshold of %s",
		formatCents(subtotal),
		formatCents(rule.thresholdCents),
	)
	if memory.PaymentReview.Status == engine.PaymentReviewRejected {
		severity = "payment-review-rejected"
		message = fmt.Sprintf(
			"Payment review was rejected for quote subtotal of %s",
			formatCents(subtotal),
		)
	}

	memory.AddFinding(engine.Finding{
		RuleName: rule.Name(),
		Severity: severity,
		Message:  message,
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
