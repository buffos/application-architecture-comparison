package products

// Reader is the public read contract provided by Products.
type Reader interface {
	GetProduct(query GetProductQuery) (ProductDetails, error)
	ListProducts(query ListProductsQuery) []ProductSummary
}

type GetProductQuery struct{ SKU string }

type ListProductsQuery struct {
	Category string
	Active   *bool
}

type ProductDetails struct {
	SKU              string
	Name             string
	Category         string
	Active           bool
	UnitPrice        int
	ReturnWindowDays int
}

type ProductSummary struct {
	SKU       string
	Name      string
	Category  string
	Active    bool
	UnitPrice int
}

func (c *Component) GetProduct(query GetProductQuery) (ProductDetails, error) {
	product, ok := c.products[query.SKU]
	if !ok {
		return ProductDetails{}, ErrProductNotFound
	}
	return productDetails(product), nil
}

func (c *Component) ListProducts(query ListProductsQuery) []ProductSummary {
	products := make([]ProductSummary, 0, len(c.products))
	for _, product := range c.products {
		if query.Category != "" && product.Category != query.Category {
			continue
		}
		if query.Active != nil && product.Active != *query.Active {
			continue
		}
		products = append(products, ProductSummary{SKU: product.SKU, Name: product.Name, Category: product.Category, Active: product.Active, UnitPrice: product.UnitPrice})
	}
	return products
}

func productDetails(product Product) ProductDetails {
	return ProductDetails{SKU: product.SKU, Name: product.Name, Category: product.Category, Active: product.Active, UnitPrice: product.UnitPrice, ReturnWindowDays: product.ReturnWindowDays}
}
