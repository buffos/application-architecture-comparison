package inventory

type ReservationRequest struct {
	SKU      ProductSKU
	Quantity int
}

type Reservation struct {
	SKU      ProductSKU
	Quantity int
}

type InventoryReservationService struct{}

func NewInventoryReservationService() InventoryReservationService {
	return InventoryReservationService{}
}

// ReserveAll coordinates several StockRecord aggregates and rolls back
// earlier reservations if a later request cannot be satisfied.
func (InventoryReservationService) ReserveAll(records map[ProductSKU]*StockRecord, requests []ReservationRequest) ([]Reservation, error) {
	reserved := make([]Reservation, 0, len(requests))
	for _, request := range requests {
		record, ok := records[request.SKU]
		if !ok {
			return nil, rollbackReservations(records, reserved, ErrProductSKURequired)
		}
		if err := record.Reserve(request.Quantity); err != nil {
			return nil, rollbackReservations(records, reserved, err)
		}
		reserved = append(reserved, Reservation{SKU: request.SKU, Quantity: request.Quantity})
	}
	return reserved, nil
}

func rollbackReservations(records map[ProductSKU]*StockRecord, reservations []Reservation, cause error) error {
	for _, reservation := range reservations {
		_ = records[reservation.SKU].Release(reservation.Quantity)
	}
	return cause
}
