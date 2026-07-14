package orders

type Repository interface {
	FindByID(id string) (Order, error)
	ListByStatus(status string) ([]Order, error)
	Save(order Order) error
}
