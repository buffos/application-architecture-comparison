package inventory

import "microkernel-architecture/internal/kernel"

type Service struct {
	stock Repository
}

func NewService(stock Repository) Service {
	return Service{
		stock: stock,
	}
}

func (s Service) Reserve(items []kernel.InventoryReservationItem) error {
	updated := make([]StockRecord, 0, len(items))

	for _, item := range items {
		record, err := s.stock.FindBySKU(item.ProductSKU)
		if err != nil {
			return err
		}

		if record.Available < item.Quantity {
			return ErrInsufficientStock
		}

		record.Available -= item.Quantity
		updated = append(updated, record)
	}

	for _, record := range updated {
		if err := s.stock.Save(record); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) Release(items []kernel.InventoryReservationItem) error {
	return s.addStock(items)
}

func (s Service) Restock(items []kernel.InventoryReservationItem) error {
	return s.addStock(items)
}

func (s Service) ListStock() ([]kernel.StockSnapshot, error) {
	stock, err := s.stock.ListStock()
	if err != nil {
		return nil, err
	}

	snapshots := make([]kernel.StockSnapshot, 0, len(stock))
	for _, item := range stock {
		snapshots = append(snapshots, kernel.StockSnapshot{
			ProductSKU: item.ProductSKU,
			Available:  item.Available,
		})
	}

	return snapshots, nil
}

func (s Service) addStock(items []kernel.InventoryReservationItem) error {
	updated := make([]StockRecord, 0, len(items))

	for _, item := range items {
		record, err := s.stock.FindBySKU(item.ProductSKU)
		if err != nil {
			return err
		}

		record.Available += item.Quantity
		updated = append(updated, record)
	}

	for _, record := range updated {
		if err := s.stock.Save(record); err != nil {
			return err
		}
	}

	return nil
}
