package returns

type Repository interface {
	FindByID(id string) (ReturnRequest, error)
	ListByStatus(status string) ([]ReturnRequest, error)
	Save(request ReturnRequest) error
}
