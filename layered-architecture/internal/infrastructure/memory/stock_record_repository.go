package memory

import (
	"sync"

	"layered-architecture/internal/domain"
)

type StockRecordRepository struct {
	mu      sync.RWMutex
	records map[string]domain.StockRecord
}

func NewStockRecordRepository() *StockRecordRepository {
	return &StockRecordRepository{
		records: make(map[string]domain.StockRecord),
	}
}

func (r *StockRecordRepository) Save(stock domain.StockRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[stock.SKU] = stock
	return nil
}

func (r *StockRecordRepository) FindBySKU(sku string) (domain.StockRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, ok := r.records[sku]
	if !ok {
		return domain.StockRecord{}, domain.ErrStockRecordNotFound
	}

	return stock, nil
}

func (r *StockRecordRepository) List() ([]domain.StockRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := make([]domain.StockRecord, 0, len(r.records))
	for _, record := range r.records {
		records = append(records, record)
	}

	return records, nil
}
