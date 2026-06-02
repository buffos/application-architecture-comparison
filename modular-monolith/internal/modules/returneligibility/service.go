package returneligibility

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Allows(request ReviewRequest) bool {
	return request.Reason != "outside return window"
}
