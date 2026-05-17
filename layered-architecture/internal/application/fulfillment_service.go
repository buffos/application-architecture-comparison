package application

import "layered-architecture/internal/domain"

type ShipmentRepository interface {
	Save(shipment domain.Shipment) error
	FindByID(id string) (domain.Shipment, error)
}

type FulfillmentService struct {
	orderRepo    OrderRepository
	stockRepo    StockRecordRepository
	shipmentRepo ShipmentRepository
}

func NewFulfillmentService(orderRepo OrderRepository, stockRepo StockRecordRepository, shipmentRepo ShipmentRepository) FulfillmentService {
	return FulfillmentService{
		orderRepo:    orderRepo,
		stockRepo:    stockRepo,
		shipmentRepo: shipmentRepo,
	}
}

func (s FulfillmentService) CreateShipment(orderID string) (domain.Shipment, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return domain.Shipment{}, err
	}

	shipment, err := domain.NewShipment(order)
	if err != nil {
		return domain.Shipment{}, err
	}

	for _, line := range order.Lines {
		stock, err := s.stockRepo.FindBySKU(line.SKU)
		if err != nil {
			return domain.Shipment{}, err
		}

		if err := stock.ConsumeReserved(line.Quantity); err != nil {
			return domain.Shipment{}, err
		}

		if err := s.stockRepo.Save(stock); err != nil {
			return domain.Shipment{}, err
		}
	}

	if err := order.MarkShipped(); err != nil {
		return domain.Shipment{}, err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return domain.Shipment{}, err
	}

	if err := s.shipmentRepo.Save(shipment); err != nil {
		return domain.Shipment{}, err
	}

	return shipment, nil
}
