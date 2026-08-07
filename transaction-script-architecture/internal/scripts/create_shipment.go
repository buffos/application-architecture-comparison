package scripts

import (
	"errors"
	"fmt"

	"transaction-script-architecture/internal/data"
)

var (
	ErrOrderNotShippable = errors.New("order is not ready for fulfillment")
	ErrNoShipmentLines   = errors.New("order has no remaining shippable lines")
	ErrShippedByRequired = errors.New("shipper is required")
)

// CreateShipment ships every remaining reserved line for a paid order and
// consumes the reserved stock as part of the same in-memory transaction.
func CreateShipment(store *data.Store, orderID string, shippedBy string) (data.Shipment, error) {
	if store == nil {
		return data.Shipment{}, ErrStoreRequired
	}

	if orderID == "" {
		return data.Shipment{}, ErrOrderIDRequired
	}

	if shippedBy == "" {
		return data.Shipment{}, ErrShippedByRequired
	}

	order, ok := store.Orders[orderID]
	if !ok {
		return data.Shipment{}, ErrOrderNotFound
	}

	if order.Status != data.OrderStatusReadyForFulfillment {
		return data.Shipment{}, ErrOrderNotShippable
	}

	shipmentLines := make([]data.ShipmentLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		remaining := line.ReservedQuantity - line.ShippedQuantity
		if remaining <= 0 {
			continue
		}

		stock, ok := store.Stocks[line.SKU]
		if !ok || stock.Reserved < remaining || stock.OnHand < remaining {
			return data.Shipment{}, ErrInsufficientStock
		}

		shipmentLines = append(shipmentLines, data.ShipmentLine{
			OrderLineID: line.ID,
			SKU:         line.SKU,
			Quantity:    remaining,
		})
	}

	if len(shipmentLines) == 0 {
		return data.Shipment{}, ErrNoShipmentLines
	}

	store.NextShipmentNumber++
	shipment := data.Shipment{
		ID:        fmt.Sprintf("shipment-%03d", store.NextShipmentNumber),
		OrderID:   order.ID,
		Status:    data.ShipmentStatusShipped,
		ShippedBy: shippedBy,
		Lines:     shipmentLines,
	}

	for _, shipmentLine := range shipmentLines {
		for index := range order.Lines {
			if order.Lines[index].ID != shipmentLine.OrderLineID {
				continue
			}

			order.Lines[index].ShippedQuantity += shipmentLine.Quantity
			stock := store.Stocks[shipmentLine.SKU]
			stock.OnHand -= shipmentLine.Quantity
			stock.Reserved -= shipmentLine.Quantity
			store.Stocks[shipmentLine.SKU] = stock
			break
		}
	}

	order.Status = data.OrderStatusShipped
	store.Orders[order.ID] = order
	store.Shipments[shipment.ID] = shipment

	return shipment, nil
}
