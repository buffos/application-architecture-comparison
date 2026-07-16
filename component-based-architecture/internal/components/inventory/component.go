package inventory

// Component owns inventory behavior and stock state for this lesson.
type Component struct {
	stock map[string]int
}

func NewComponent() *Component {
	return &Component{stock: make(map[string]int)}
}

func (c *Component) RegisterStock(record StockRecord) {
	c.stock[record.ProductSKU] = record.Available
}

func (c *Component) Reserve(items []ReservationItem) error {
	requested := make(map[string]int, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrReservationQuantityMustBePositive
		}
		requested[item.ProductSKU] += item.Quantity
	}

	for productSKU, quantity := range requested {
		available, ok := c.stock[productSKU]
		if !ok {
			return ErrStockNotFound
		}
		if available < quantity {
			return ErrInsufficientStock
		}
	}

	for productSKU, quantity := range requested {
		c.stock[productSKU] -= quantity
	}
	return nil
}

func (c *Component) Release(items []ReleaseItem) error {
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrReservationQuantityMustBePositive
		}
		if _, ok := c.stock[item.ProductSKU]; !ok {
			return ErrStockNotFound
		}
	}
	for _, item := range items {
		c.stock[item.ProductSKU] += item.Quantity
	}
	return nil
}
func (c *Component) Restock(items []RestockItem) error {
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrReservationQuantityMustBePositive
		}
		if _, ok := c.stock[item.ProductSKU]; !ok {
			return ErrStockNotFound
		}
	}
	for _, item := range items {
		c.stock[item.ProductSKU] += item.Quantity
	}
	return nil
}

var _ StockKeeper = (*Component)(nil)
