package returneligibility

// Component owns return-acceptance policy for this lesson.
type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Allows(review Review) bool {
	if review.ShippedAt.IsZero() || review.RequestedAt.Before(review.ShippedAt) {
		return false
	}
	for _, line := range review.Lines {
		if review.RequestedAt.After(review.ShippedAt.AddDate(0, 0, line.ReturnWindowDays)) {
			return false
		}
	}
	return true
}

var _ Evaluator = (*Component)(nil)
