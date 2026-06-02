package returns

import "time"

type ReturnableOrder struct {
	OrderID    string
	CustomerID string
	ShippedAt  time.Time
	Lines      []ReturnableOrderLine
}

type ReturnableOrderLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	ShippedQuantity  int
	UnitPrice        int
	ReturnWindowDays int
}
