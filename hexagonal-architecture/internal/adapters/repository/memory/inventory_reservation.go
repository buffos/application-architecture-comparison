package memory

import "hexagonal-architecture/internal/core/domain"

type InventoryReservationAdapter struct {
	available map[string]int
}

func NewInventoryReservationAdapter(initial map[string]int) *InventoryReservationAdapter {
	copyMap := make(map[string]int, len(initial))
	for sku, quantity := range initial {
		copyMap[sku] = quantity
	}

	return &InventoryReservationAdapter{available: copyMap}
}

func (a *InventoryReservationAdapter) Reserve(lines []domain.ReservationLine) error {
	for _, line := range lines {
		if a.available[line.SKU] < line.Quantity {
			return domain.ErrInsufficientStock
		}
	}

	for _, line := range lines {
		a.available[line.SKU] -= line.Quantity
	}

	return nil
}

func (a *InventoryReservationAdapter) Available(sku string) int {
	return a.available[sku]
}
