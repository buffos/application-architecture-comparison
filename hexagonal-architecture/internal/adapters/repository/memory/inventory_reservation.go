package memory

import "hexagonal-architecture/internal/core/domain"

type InventoryReservationAdapter struct {
	available  map[string]int
	thresholds map[string]int
}

func NewInventoryReservationAdapter(initial map[string]int) *InventoryReservationAdapter {
	copyMap := make(map[string]int, len(initial))
	for sku, quantity := range initial {
		copyMap[sku] = quantity
	}

	return &InventoryReservationAdapter{
		available:  copyMap,
		thresholds: make(map[string]int, len(initial)),
	}
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

func (a *InventoryReservationAdapter) Consume(lines []domain.ReservationLine) error {
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

func (a *InventoryReservationAdapter) Release(lines []domain.ReservationLine) error {
	for _, line := range lines {
		a.available[line.SKU] += line.Quantity
	}

	return nil
}

func (a *InventoryReservationAdapter) Restock(lines []domain.ReservationLine) error {
	for _, line := range lines {
		a.available[line.SKU] += line.Quantity
	}

	return nil
}

func (a *InventoryReservationAdapter) Available(sku string) int {
	return a.available[sku]
}

func (a *InventoryReservationAdapter) SetReorderThreshold(sku string, threshold int) {
	a.thresholds[sku] = threshold
}

func (a *InventoryReservationAdapter) FindBySKU(sku string) (domain.StockRecord, error) {
	available, ok := a.available[sku]
	if !ok {
		return domain.StockRecord{}, domain.ErrStockRecordNotFound
	}

	return domain.StockRecord{
		SKU:              sku,
		Available:        available,
		ReorderThreshold: a.thresholds[sku],
	}, nil
}

func (a *InventoryReservationAdapter) Save(record domain.StockRecord) error {
	a.available[record.SKU] = record.Available
	a.thresholds[record.SKU] = record.ReorderThreshold
	return nil
}

func (a *InventoryReservationAdapter) ListStock() ([]domain.StockRecord, error) {
	records := make([]domain.StockRecord, 0, len(a.available))
	for sku, available := range a.available {
		records = append(records, domain.StockRecord{
			SKU:              sku,
			Available:        available,
			ReorderThreshold: a.thresholds[sku],
		})
	}

	return records, nil
}
