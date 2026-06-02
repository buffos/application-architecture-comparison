package returns

type ReturnableOrder struct {
	OrderID    string
	CustomerID string
	Lines      []ReturnableOrderLine
}

type ReturnableOrderLine struct {
	ProductSKU      string
	ProductName     string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}
