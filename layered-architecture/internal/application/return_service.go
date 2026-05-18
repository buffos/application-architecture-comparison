package application

import "layered-architecture/internal/domain"

type ReturnRequestRepository interface {
	Save(request domain.ReturnRequest) error
	FindByID(id string) (domain.ReturnRequest, error)
}

type ReturnService struct {
	orderRepo  OrderRepository
	stockRepo  StockRecordRepository
	returnRepo ReturnRequestRepository
}

func NewReturnService(orderRepo OrderRepository, stockRepo StockRecordRepository, returnRepo ReturnRequestRepository) ReturnService {
	return ReturnService{
		orderRepo:  orderRepo,
		stockRepo:  stockRepo,
		returnRepo: returnRepo,
	}
}

func (s ReturnService) RequestReturn(orderID string, reason string) (domain.ReturnRequest, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	request, err := domain.NewReturnRequest(order, reason)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := s.returnRepo.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}

func (s ReturnService) AcceptReturn(id string) (domain.ReturnRequest, error) {
	request, err := s.returnRepo.FindByID(id)
	if err != nil {
		return domain.ReturnRequest{}, err
	}

	if err := request.Accept(); err != nil {
		return domain.ReturnRequest{}, err
	}

	for _, line := range request.Lines {
		stock, err := s.stockRepo.FindBySKU(line.SKU)
		if err != nil {
			return domain.ReturnRequest{}, err
		}

		if err := stock.Restock(line.Quantity); err != nil {
			return domain.ReturnRequest{}, err
		}

		if err := s.stockRepo.Save(stock); err != nil {
			return domain.ReturnRequest{}, err
		}
	}

	if err := s.returnRepo.Save(request); err != nil {
		return domain.ReturnRequest{}, err
	}

	return request, nil
}
