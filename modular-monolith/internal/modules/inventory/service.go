package inventory

type Reserver interface {
	Reserve(items []ReservationItem) error
}

type Service struct {
	stock Repository
}

func NewService(stock Repository) Service {
	return Service{
		stock: stock,
	}
}

func (s Service) Reserve(items []ReservationItem) error {
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrReservationQuantityMustBePositive
		}
	}

	return s.stock.Reserve(items)
}
