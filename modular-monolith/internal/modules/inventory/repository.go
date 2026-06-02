package inventory

type Repository interface {
	Save(record StockRecord) error
	Reserve(items []ReservationItem) error
	Release(items []ReleaseItem) error
}
