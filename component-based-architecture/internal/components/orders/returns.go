package orders

import "time"

type ReturnableOrderSource interface {
	GetReturnableOrder(orderID string) (ReturnableOrder, error)
}
type ReturnableOrder struct {
	OrderID, CustomerID string
	Lines               []ReturnableOrderLine
	ShippedAt           time.Time
}
type ReturnableOrderLine struct {
	ProductSKU          string
	Quantity, UnitPrice int
	ReturnWindowDays    int
}

func (c *Component) GetReturnableOrder(orderID string) (ReturnableOrder, error) {
	order, ok := c.orders[orderID]
	if !ok {
		return ReturnableOrder{}, ErrOrderNotFound
	}
	if err := order.EnsureReturnable(); err != nil {
		return ReturnableOrder{}, err
	}
	lines := make([]ReturnableOrderLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ReturnableOrderLine{ProductSKU: line.ProductSKU, Quantity: line.Quantity, UnitPrice: line.UnitPrice, ReturnWindowDays: line.ReturnWindowDays})
	}
	return ReturnableOrder{OrderID: order.ID, CustomerID: order.CustomerID, Lines: lines, ShippedAt: order.ShippedAt}, nil
}

var _ ReturnableOrderSource = (*Component)(nil)
