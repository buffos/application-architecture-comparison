package shipments

type Repository interface {
	FindByID(id string) (Shipment, error)
	Save(shipment Shipment) error
}
