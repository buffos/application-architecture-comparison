package customers

// CustomerDirectory is the public contract that this component provides to
// other components. It deliberately exposes no customer storage details.
type CustomerDirectory interface {
	RequireActiveCustomer(id string) error
}
