package quotes

type Repository interface {
	Save(quote Quote) error
	FindByID(id string) (Quote, error)
	ListByStatus(status string) ([]Quote, error)
}
