package shipments

import (
	"errors"
	"sort"

	"domain-driven-design-architecture/internal/domain/fulfillment"
)

var ErrShipmentNotFound = errors.New("shipment not found")

type Reader interface {
	GetShipment(id string) (ShipmentDetails, error)
	ListShipments(status fulfillment.ShipmentStatus) []ShipmentSummary
}

type ShipmentDetails struct {
	ShipmentID string
	OrderID    string
	Status     string
	Lines      []ShipmentLineDetails
}

type ShipmentLineDetails struct {
	ProductSKU string
	Quantity   int
}

type ShipmentSummary struct {
	ShipmentID string
	OrderID    string
	Status     string
	LineCount  int
}

type InMemoryReader struct {
	shipments map[string]ShipmentDetails
}

func NewInMemoryReader() *InMemoryReader {
	return &InMemoryReader{shipments: make(map[string]ShipmentDetails)}
}

func (r *InMemoryReader) Save(shipment fulfillment.Shipment) {
	details := ShipmentDetails{ShipmentID: string(shipment.ID()), OrderID: string(shipment.OrderID()), Status: string(shipment.Status()), Lines: make([]ShipmentLineDetails, 0, len(shipment.Lines()))}
	for _, line := range shipment.Lines() {
		details.Lines = append(details.Lines, ShipmentLineDetails{ProductSKU: string(line.ProductSKU()), Quantity: line.Quantity()})
	}
	r.shipments[details.ShipmentID] = details
}

func (r *InMemoryReader) GetShipment(id string) (ShipmentDetails, error) {
	details, ok := r.shipments[id]
	if !ok {
		return ShipmentDetails{}, ErrShipmentNotFound
	}
	details.Lines = append([]ShipmentLineDetails(nil), details.Lines...)
	return details, nil
}

func (r *InMemoryReader) ListShipments(status fulfillment.ShipmentStatus) []ShipmentSummary {
	result := make([]ShipmentSummary, 0, len(r.shipments))
	for _, details := range r.shipments {
		if status != "" && details.Status != string(status) {
			continue
		}
		result = append(result, ShipmentSummary{ShipmentID: details.ShipmentID, OrderID: details.OrderID, Status: details.Status, LineCount: len(details.Lines)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ShipmentID < result[j].ShipmentID })
	return result
}

var _ Reader = (*InMemoryReader)(nil)
