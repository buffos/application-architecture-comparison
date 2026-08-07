package inventory

type RestockRequest struct {
	SKU      ProductSKU
	Quantity int
}

type ReturnRestockingService struct{}

func NewReturnRestockingService() ReturnRestockingService {
	return ReturnRestockingService{}
}

func (ReturnRestockingService) RestockAll(records map[ProductSKU]*StockRecord, requests []RestockRequest) error {
	for _, request := range requests {
		record, ok := records[request.SKU]
		if !ok {
			return ErrProductSKURequired
		}
		if err := record.Receive(request.Quantity); err != nil {
			return err
		}
	}
	return nil
}
