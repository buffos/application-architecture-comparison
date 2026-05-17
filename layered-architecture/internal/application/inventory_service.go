package application

import "layered-architecture/internal/domain"

type StockRecordRepository interface {
	Save(stock domain.StockRecord) error
	FindBySKU(sku string) (domain.StockRecord, error)
}

type InventoryService struct {
	productRepo ProductRepository
	stockRepo   StockRecordRepository
}

func NewInventoryService(productRepo ProductRepository, stockRepo StockRecordRepository) InventoryService {
	return InventoryService{
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

func (s InventoryService) ReceiveStock(sku string, quantity int) (domain.StockRecord, error) {
	product, err := s.productRepo.FindBySKU(sku)
	if err != nil {
		return domain.StockRecord{}, err
	}

	stock, err := s.stockRepo.FindBySKU(product.SKU)
	if err != nil {
		if !isNotFound(err, domain.ErrStockRecordNotFound) {
			return domain.StockRecord{}, err
		}

		stock, err = domain.NewStockRecord(product.SKU, domain.StockShortageRejectOrder)
		if err != nil {
			return domain.StockRecord{}, err
		}
	}

	if err := stock.Receive(quantity); err != nil {
		return domain.StockRecord{}, err
	}

	if err := s.stockRepo.Save(stock); err != nil {
		return domain.StockRecord{}, err
	}

	return stock, nil
}

func isNotFound(err error, target error) bool {
	return err == target
}
