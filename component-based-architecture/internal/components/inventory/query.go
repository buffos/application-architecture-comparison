package inventory

// StockReader is the narrow read contract provided to operational reports.
type StockReader interface {
	ListStock() []StockSnapshot
}

type StockSnapshot struct {
	ProductSKU string
	Available  int
}

func (c *Component) ListStock() []StockSnapshot {
	stock := make([]StockSnapshot, 0, len(c.stock))
	for sku, available := range c.stock {
		stock = append(stock, StockSnapshot{ProductSKU: sku, Available: available})
	}
	return stock
}
