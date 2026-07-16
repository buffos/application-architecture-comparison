package products

// Component owns product behavior and its in-memory state for this lesson.
type Component struct {
	products map[string]Product
}

func NewComponent() *Component {
	return &Component{products: make(map[string]Product)}
}

func (c *Component) Register(product Product) error {
	if product.SKU == "" {
		return ErrProductSKURequired
	}
	c.products[product.SKU] = product
	return nil
}

func (c *Component) GetProductForQuote(sku string) (ProductForQuote, error) {
	product, ok := c.products[sku]
	if !ok {
		return ProductForQuote{}, ErrProductNotFound
	}
	if !product.Active {
		return ProductForQuote{}, ErrProductInactive
	}
	return ProductForQuote{SKU: product.SKU, Name: product.Name, Category: product.Category, UnitPrice: product.UnitPrice, ReturnWindowDays: product.ReturnWindowDays}, nil
}

var _ Catalog = (*Component)(nil)
