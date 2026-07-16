package payments

// Component owns payment-processing behavior for this lesson.
type Component struct {
	gateway Gateway
}

func NewComponent(gateway Gateway) *Component {
	return &Component{gateway: gateway}
}

func (c *Component) Capture(request PaymentRequest) (CaptureResult, error) {
	return c.gateway.Capture(request)
}

var _ Processor = (*Component)(nil)
