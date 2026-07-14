package orders

type Repository interface {
	FindByID(id string) (Order, error)
	Save(order Order) error
}
