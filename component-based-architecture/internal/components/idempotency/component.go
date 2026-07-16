package idempotency

// Component owns completed command results for this in-memory lesson.
type Component struct {
	results map[string]Result
}

func NewComponent() *Component {
	return &Component{results: make(map[string]Result)}
}

func (c *Component) Find(key string) (Result, bool) {
	result, ok := c.results[key]
	return result, ok
}

func (c *Component) Save(key string, result Result) {
	c.results[key] = result
}

var _ Store = (*Component)(nil)
