package scripts

import (
	"errors"
	"fmt"
	"time"

	"transaction-script-architecture/internal/data"
)

var (
	ErrOrderNotShippable    = errors.New("order is not ready for fulfillment")
	ErrNoShipmentLines      = errors.New("order has no remaining shippable lines")
	ErrShipmentLinesInvalid = errors.New("shipment lines are invalid")
	ErrShippedByRequired    = errors.New("shipper is required")
)

// CreateShipment ships every remaining reserved line by delegating to the
// line-aware transaction.
func CreateShipment(store *data.Store, orderID string, shippedBy string) (data.Shipment, error) {
	return CreatePartialShipment(store, orderID, nil, shippedBy)
}

// CreatePartialShipment ships only the requested reserved quantities and
// derives the order's partial/full status after the write.
func CreatePartialShipment(store *data.Store, orderID string, lines []data.ShipmentLine, shippedBy string) (data.Shipment, error) {
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

	if order.Status != data.OrderStatusReadyForFulfillment && order.Status != data.OrderStatusPartiallyShipped {
		return data.Shipment{}, ErrOrderNotShippable
	}

	shipmentLines := lines
	if lines == nil {
		shipmentLines = make([]data.ShipmentLine, 0, len(order.Lines))
		for _, line := range order.Lines {
			remaining := line.ReservedQuantity - line.ShippedQuantity
			if remaining <= 0 {
				continue
			}
			shipmentLines = append(shipmentLines, data.ShipmentLine{
				OrderLineID: line.ID,
				SKU:         line.SKU,
				Quantity:    remaining,
			})
		}
	}

	if len(shipmentLines) == 0 {
		return data.Shipment{}, ErrNoShipmentLines
	}

	for _, shipmentLine := range shipmentLines {
		if shipmentLine.Quantity <= 0 {
			return data.Shipment{}, ErrShipmentLinesInvalid
		}

		matched := false
		for _, orderLine := range order.Lines {
			if orderLine.ID != shipmentLine.OrderLineID {
				continue
			}

			remaining := orderLine.ReservedQuantity - orderLine.ShippedQuantity
			if shipmentLine.SKU != orderLine.SKU || shipmentLine.Quantity > remaining {
				return data.Shipment{}, ErrShipmentLinesInvalid
			}

			stock, stockExists := store.Stocks[shipmentLine.SKU]
			if !stockExists || stock.Reserved < shipmentLine.Quantity || stock.OnHand < shipmentLine.Quantity {
				return data.Shipment{}, ErrInsufficientStock
			}

			matched = true
			break
		}

		if !matched {
			return data.Shipment{}, ErrShipmentLinesInvalid
		}
	}

	store.NextShipmentNumber++
	shipment := data.Shipment{
		ID:        fmt.Sprintf("shipment-%03d", store.NextShipmentNumber),
		OrderID:   order.ID,
		Status:    data.ShipmentStatusShipped,
		ShippedBy: shippedBy,
		ShippedAt: time.Now(),
		Lines:     append([]data.ShipmentLine(nil), shipmentLines...),
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

	allShipped := true
	for _, line := range order.Lines {
		if line.ShippedQuantity < line.OrderedQuantity {
			allShipped = false
			break
		}
	}
	if allShipped {
		order.Status = data.OrderStatusShipped
	} else {
		order.Status = data.OrderStatusPartiallyShipped
	}
	order.ShippedAt = shipment.ShippedAt
	store.Orders[order.ID] = order
	store.Shipments[shipment.ID] = shipment

	return shipment, nil
}
