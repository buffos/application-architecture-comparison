package shipments

type Repository interface {
	Save(shipment Shipment) error
	FindByID(id string) (Shipment, error)
}
