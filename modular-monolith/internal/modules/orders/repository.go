package orders

type Repository interface {
	Save(order Order) error
	FindByID(id string) (Order, error)
	ListByStatus(status string) ([]Order, error)
}
