package quotes

type Repository interface {
	FindByID(id string) (Quote, error)
	ListByStatus(status string) ([]Quote, error)
	Save(quote Quote) error
}
