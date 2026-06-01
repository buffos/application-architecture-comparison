package application

import "onion-architecture/internal/domain"

type CancelOrderCommand struct {
	OrderID string
}

type CancelOrderResult struct {
	OrderID string
	Status  string
}

type InventoryRelease interface {
	Release(items []domain.InventoryReleaseItem) error
}

type CancelOrderService struct {
	orders    OrderRepository
	inventory InventoryRelease
}

func NewCancelOrderService(orders OrderRepository, inventory InventoryRelease) CancelOrderService {
	return CancelOrderService{
		orders:    orders,
		inventory: inventory,
	}
}

func (s CancelOrderService) Execute(command CancelOrderCommand) (CancelOrderResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return CancelOrderResult{}, err
	}

	if err := order.Cancel(); err != nil {
		return CancelOrderResult{}, err
	}

	items := make([]domain.InventoryReleaseItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		items = append(items, domain.InventoryReleaseItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.inventory.Release(items); err != nil {
		return CancelOrderResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return CancelOrderResult{}, err
	}

	return CancelOrderResult{
		OrderID: order.ID,
		Status:  order.Status,
	}, nil
}
