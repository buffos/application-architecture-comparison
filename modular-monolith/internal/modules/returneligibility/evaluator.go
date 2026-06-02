package returneligibility

import "time"

type ReviewRequest struct {
	Reason      string
	ShippedAt   time.Time
	RequestedAt time.Time
	Lines       []ReviewLine
}

type ReviewLine struct {
	ReturnWindowDays int
}

type Evaluator interface {
	Allows(request ReviewRequest) bool
}
