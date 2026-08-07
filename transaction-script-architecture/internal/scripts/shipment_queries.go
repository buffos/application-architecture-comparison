package scripts

import (
	"errors"
	"sort"

	"transaction-script-architecture/internal/data"
)

var (
	ErrShipmentIDRequired = errors.New("shipment id is required")
	ErrShipmentNotFound   = errors.New("shipment not found")
)

func GetShipment(store *data.Store, shipmentID string) (data.Shipment, error) {
	if store == nil {
		return data.Shipment{}, ErrStoreRequired
	}
	if shipmentID == "" {
		return data.Shipment{}, ErrShipmentIDRequired
	}

	shipment, ok := store.Shipments[shipmentID]
	if !ok {
		return data.Shipment{}, ErrShipmentNotFound
	}

	return cloneShipment(shipment), nil
}

func ListShipments(store *data.Store, orderID string) ([]data.Shipment, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	shipments := make([]data.Shipment, 0, len(store.Shipments))
	for _, shipment := range store.Shipments {
		if orderID != "" && shipment.OrderID != orderID {
			continue
		}
		shipments = append(shipments, cloneShipment(shipment))
	}

	sort.Slice(shipments, func(i, j int) bool {
		return shipments[i].ID < shipments[j].ID
	})

	return shipments, nil
}

func cloneShipment(shipment data.Shipment) data.Shipment {
	shipment.Lines = append([]data.ShipmentLine(nil), shipment.Lines...)
	return shipment
}
