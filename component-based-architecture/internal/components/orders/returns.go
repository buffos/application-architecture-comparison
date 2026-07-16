package orders

type ReturnableOrderSource interface {
	GetReturnableOrder(orderID string) (ReturnableOrder, error)
}
type ReturnableOrder struct {
	OrderID, CustomerID string
	Lines               []ReturnableOrderLine
}
type ReturnableOrderLine struct {
	ProductSKU          string
	Quantity, UnitPrice int
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
		lines = append(lines, ReturnableOrderLine{ProductSKU: line.ProductSKU, Quantity: line.Quantity, UnitPrice: line.UnitPrice})
	}
	return ReturnableOrder{OrderID: order.ID, CustomerID: order.CustomerID, Lines: lines}, nil
}

var _ ReturnableOrderSource = (*Component)(nil)
