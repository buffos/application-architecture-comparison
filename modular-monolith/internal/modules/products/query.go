package products

type GetProductQuery struct {
	SKU string
}

type ProductDetails struct {
	SKU              string
	Name             string
	Category         string
	Active           bool
	UnitPrice        int
	ReturnWindowDays int
}

type ListProductsQuery struct {
	Category   string
	ActiveOnly bool
}

func (s Service) GetProduct(query GetProductQuery) (ProductDetails, error) {
	product, err := s.products.FindBySKU(query.SKU)
	if err != nil {
		return ProductDetails{}, err
	}

	return ProductDetails{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		Active:           product.Active,
		UnitPrice:        product.UnitPrice,
		ReturnWindowDays: product.ReturnWindowDays,
	}, nil
}

func (s Service) ListProducts(query ListProductsQuery) ([]ProductDetails, error) {
	products, err := s.products.List(query.Category, query.ActiveOnly)
	if err != nil {
		return nil, err
	}

	list := make([]ProductDetails, 0, len(products))
	for _, product := range products {
		list = append(list, ProductDetails{
			SKU:              product.SKU,
			Name:             product.Name,
			Category:         product.Category,
			Active:           product.Active,
			UnitPrice:        product.UnitPrice,
			ReturnWindowDays: product.ReturnWindowDays,
		})
	}

	return list, nil
}
