package customers

type Directory interface {
	RequireActiveCustomer(id string) error
}
