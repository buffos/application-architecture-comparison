package orders

import (
	"errors"
	"sort"

	"domain-driven-design-architecture/internal/domain/ordering"
)

var ErrOrderNotFound = errors.New("order not found")

type Reader interface {
	GetOrder(id string) (OrderDetails, error)
	ListOrders(status ordering.OrderStatus) []OrderSummary
}

type OrderDetails struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	TotalCents int64
	Currency   string
	Lines      []OrderLineDetails
}

type OrderLineDetails struct {
	ProductSKU string
	Category   string
	Quantity   int
}

type OrderSummary struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type InMemoryReader struct {
	orders map[string]OrderDetails
}

func NewInMemoryReader() *InMemoryReader {
	return &InMemoryReader{orders: make(map[string]OrderDetails)}
}

func (r *InMemoryReader) Save(order ordering.Order) error {
	total, err := order.Total()
	if err != nil {
		return err
	}
	details := OrderDetails{OrderID: string(order.ID()), QuoteID: string(order.QuoteID()), CustomerID: string(order.CustomerID()), Status: string(order.Status()), TotalCents: total.Cents(), Currency: total.Currency(), Lines: make([]OrderLineDetails, 0, len(order.Lines()))}
	for _, line := range order.Lines() {
		details.Lines = append(details.Lines, OrderLineDetails{ProductSKU: string(line.ProductSKU()), Category: string(line.ProductCategory()), Quantity: line.Quantity()})
	}
	r.orders[details.OrderID] = details
	return nil
}

func (r *InMemoryReader) GetOrder(id string) (OrderDetails, error) {
	details, ok := r.orders[id]
	if !ok {
		return OrderDetails{}, ErrOrderNotFound
	}
	details.Lines = append([]OrderLineDetails(nil), details.Lines...)
	return details, nil
}

func (r *InMemoryReader) ListOrders(status ordering.OrderStatus) []OrderSummary {
	result := make([]OrderSummary, 0, len(r.orders))
	for _, details := range r.orders {
		if status != "" && details.Status != string(status) {
			continue
		}
		result = append(result, OrderSummary{OrderID: details.OrderID, CustomerID: details.CustomerID, Status: details.Status, LineCount: len(details.Lines)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].OrderID < result[j].OrderID })
	return result
}

var _ Reader = (*InMemoryReader)(nil)
