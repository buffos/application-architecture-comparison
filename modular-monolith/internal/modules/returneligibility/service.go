package returneligibility

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Allows(request ReviewRequest) bool {
	if request.Reason == "outside return window" {
		return false
	}

	if request.ShippedAt.IsZero() || request.RequestedAt.Before(request.ShippedAt) {
		return false
	}

	for _, line := range request.Lines {
		deadline := request.ShippedAt.AddDate(0, 0, line.ReturnWindowDays)
		if request.RequestedAt.After(deadline) {
			return false
		}
	}

	return true
}
