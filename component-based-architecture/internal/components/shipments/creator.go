package shipments

// Creator is the public contract this component provides to fulfillment
// workflows that need a shipment record.
type Creator interface {
	Create(request ShipmentRequest) (Shipment, error)
}
