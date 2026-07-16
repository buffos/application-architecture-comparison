package products

// Catalog is the public contract that this component provides to consumers
// that need a sellable product snapshot.
type Catalog interface {
	GetProductForQuote(sku string) (ProductForQuote, error)
}
