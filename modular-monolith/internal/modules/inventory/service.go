package inventory

type Reserver interface {
	Reserve(items []ReservationItem) error
}

type Releaser interface {
	Release(items []ReleaseItem) error
}

type StockKeeper interface {
	Reserver
	Releaser
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

func (s Service) Release(items []ReleaseItem) error {
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrReservationQuantityMustBePositive
		}
	}

	return s.stock.Release(items)
}
