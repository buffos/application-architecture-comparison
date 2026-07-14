package payments

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Capture(orderID string, amount int) error {
	return nil
}
