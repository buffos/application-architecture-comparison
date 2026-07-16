package returneligibility

import "time"

// Evaluator is the public policy contract provided by this component to
// workflows that need a return-acceptance decision.
type Evaluator interface {
	Allows(review Review) bool
}

// Review is the small policy snapshot Returns shares without exposing its
// private return-request storage.
type Review struct {
	ShippedAt   time.Time
	RequestedAt time.Time
	Lines       []ReviewLine
}

type ReviewLine struct {
	ReturnWindowDays int
}
