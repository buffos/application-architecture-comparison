package orders

import "component-based-architecture/internal/components/quotes"

const OrderStatusPendingPayment = "PendingPayment"

type Order struct {
	ID         string
	QuoteID    string
	CustomerID string
	Status     string
	Lines      []OrderLine
}

type OrderLine struct {
	ProductSKU      string
	ProductName     string
	ProductCategory string
	Quantity        int
	UnitPrice       int
}

func newOrderFromApprovedQuote(id string, quote quotes.ApprovedQuote) Order {
	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU: line.ProductSKU, ProductName: line.ProductName, ProductCategory: line.ProductCategory,
			Quantity: line.Quantity, UnitPrice: line.UnitPrice,
		})
	}
	return Order{
		ID: id, QuoteID: quote.QuoteID, CustomerID: quote.CustomerID, Status: OrderStatusPendingPayment, Lines: lines,
	}
}
