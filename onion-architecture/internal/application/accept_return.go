package application

import "onion-architecture/internal/domain"

type RefundGateway interface {
	Refund(order domain.Order) error
}

type InventoryRestock interface {
	Restock(items []domain.InventoryRestockItem) error
}

type ReturnEligibilityPolicy interface {
	IsEligible(request domain.ReturnRequest, order domain.Order) (bool, error)
}

type AcceptReturnCommand struct {
	ReturnRequestID string
	IdempotencyKey  string
	ReviewedBy      string
	ReviewNote      string
	ProcessedBy     string
}

type AcceptReturnResult struct {
	ReturnRequestID string
	Status          string
}

type AcceptReturnService struct {
	orders  OrderRepository
	returns ReturnRequestStore
	policy  ReturnEligibilityPolicy
	idempotency IdempotencyStore
	refunds RefundGateway
	restock InventoryRestock
}

func NewAcceptReturnService(orders OrderRepository, returns ReturnRequestStore, policy ReturnEligibilityPolicy, idempotency IdempotencyStore, refunds RefundGateway, restock InventoryRestock) AcceptReturnService {
	return AcceptReturnService{
		orders:  orders,
		returns: returns,
		policy:  policy,
		idempotency: idempotency,
		refunds: refunds,
		restock: restock,
	}
}

func (s AcceptReturnService) Execute(command AcceptReturnCommand) (AcceptReturnResult, error) {
	if status, ok, err := s.idempotency.Get(command.IdempotencyKey); err != nil {
		return AcceptReturnResult{}, err
	} else if ok {
		return AcceptReturnResult{
			ReturnRequestID: command.ReturnRequestID,
			Status:          status,
		}, nil
	}

	request, err := s.returns.FindByID(command.ReturnRequestID)
	if err != nil {
		return AcceptReturnResult{}, err
	}

	order, err := s.orders.FindByID(request.OrderID)
	if err != nil {
		return AcceptReturnResult{}, err
	}

	eligible, err := s.policy.IsEligible(request, order)
	if err != nil {
		return AcceptReturnResult{}, err
	}

	if !eligible {
		if err := s.idempotency.Save(command.IdempotencyKey, request.Status); err != nil {
			return AcceptReturnResult{}, err
		}

		return AcceptReturnResult{
			ReturnRequestID: request.ID,
			Status:          request.Status,
		}, nil
	}

	if err := request.Accept(command.ReviewedBy, command.ReviewNote); err != nil {
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

	if err := request.MarkRefunded(command.ProcessedBy); err != nil {
		return AcceptReturnResult{}, err
	}

	if err := s.returns.Save(request); err != nil {
		return AcceptReturnResult{}, err
	}

	if err := s.idempotency.Save(command.IdempotencyKey, request.Status); err != nil {
		return AcceptReturnResult{}, err
	}

	return AcceptReturnResult{
		ReturnRequestID: request.ID,
		Status:          request.Status,
	}, nil
}
