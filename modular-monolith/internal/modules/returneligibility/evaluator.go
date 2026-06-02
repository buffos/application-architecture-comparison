package returneligibility

type ReviewRequest struct {
	Reason string
}

type Evaluator interface {
	Allows(request ReviewRequest) bool
}
