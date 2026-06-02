package inventory

type Reserver interface {
	Reserve(items []ReservationItem) error
}

type Releaser interface {
	Release(items []ReleaseItem) error
}

type Restocker interface {
	Restock(items []RestockItem) error
}

type StockKeeper interface {
	Reserver
	Releaser
}

type StockSnapshot struct {
	ProductSKU string
	Available  int
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

func (s Service) Restock(items []RestockItem) error {
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrReservationQuantityMustBePositive
		}
	}

	return s.stock.Restock(items)
}

func (s Service) ListStock() ([]StockSnapshot, error) {
	records, err := s.stock.List()
	if err != nil {
		return nil, err
	}

	list := make([]StockSnapshot, 0, len(records))
	for _, record := range records {
		list = append(list, StockSnapshot{
			ProductSKU: record.ProductSKU,
			Available:  record.Available,
		})
	}

	return list, nil
}
