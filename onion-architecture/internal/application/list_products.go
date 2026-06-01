package application

type ListProductsQuery struct {
	Category   string
	ActiveOnly bool
}

type ListProductsService struct {
	products ProductFinder
}

func NewListProductsService(products ProductFinder) ListProductsService {
	return ListProductsService{
		products: products,
	}
}

func (s ListProductsService) Execute(query ListProductsQuery) ([]ProductDetails, error) {
	products, err := s.products.List(query.Category, query.ActiveOnly)
	if err != nil {
		return nil, err
	}

	result := make([]ProductDetails, 0, len(products))
	for _, product := range products {
		result = append(result, ProductDetails{
			SKU:              product.SKU,
			Name:             product.Name,
			Category:         product.Category,
			Active:           product.Active,
			UnitPrice:        product.UnitPrice,
			ReturnWindowDays: product.ReturnWindowDays,
		})
	}

	return result, nil
}
