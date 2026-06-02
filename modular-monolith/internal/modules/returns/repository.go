package returns

type Repository interface {
	Save(request ReturnRequest) error
	FindByID(id string) (ReturnRequest, error)
	ListByStatus(status string) ([]ReturnRequest, error)
}
