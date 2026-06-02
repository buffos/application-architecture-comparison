package products

type Catalog interface {
	GetProductForQuote(sku string) (ProductForQuote, error)
}
