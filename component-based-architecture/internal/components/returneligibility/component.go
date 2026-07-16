package returneligibility

// Component owns return-acceptance policy for this lesson.
type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Allows(review Review) bool {
	return review.Reason != "outside return window"
}

var _ Evaluator = (*Component)(nil)
