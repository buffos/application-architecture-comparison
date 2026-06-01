package application

import "onion-architecture/internal/domain"

type RefundGateway interface {
	Refund(order domain.Order) error
}

type InventoryRestock interface {
	Restock(items []domain.InventoryRestockItem) error
}

type AcceptReturnCommand struct {
	ReturnRequestID string
}

type AcceptReturnResult struct {
	ReturnRequestID string
	Status          string
}

type AcceptReturnService struct {
	orders  OrderRepository
	returns ReturnRequestStore
	refunds RefundGateway
	restock InventoryRestock
}

func NewAcceptReturnService(orders OrderRepository, returns ReturnRequestStore, refunds RefundGateway, restock InventoryRestock) AcceptReturnService {
	return AcceptReturnService{
		orders:  orders,
		returns: returns,
		refunds: refunds,
		restock: restock,
	}
}

func (s AcceptReturnService) Execute(command AcceptReturnCommand) (AcceptReturnResult, error) {
	request, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return AcceptReturnResult{}, err
	}

	order, err := s.orders.FindByID(request.OrderID)
	if err != nil {
		return AcceptReturnResult{}, err
	}

	if err := request.Accept(); err != nil {
		return AcceptReturnResult{}, err
	}

	if err := s.refunds.Refund(order); err != nil {
		return AcceptReturnResult{}, err
	}

	items := make([]domain.InventoryRestockItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		items = append(items, domain.InventoryRestockItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.restock.Restock(items); err != nil {
		return AcceptReturnResult{}, err
	}

	if err := request.MarkRefunded(); err != nil {
		return AcceptReturnResult{}, err
	}

	if err := s.returns.Save(request); err != nil {
		return AcceptReturnResult{}, err
	}

	return AcceptReturnResult{
		ReturnRequestID: request.ID,
		Status:          request.Status,
	}, nil
}
