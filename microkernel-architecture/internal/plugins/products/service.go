package products

import "microkernel-architecture/internal/kernel"

type Service struct {
	products Repository
}

func NewService(products Repository) Service {
	return Service{
		products: products,
	}
}

func (s Service) GetProductForQuote(sku string) (kernel.Product, error) {
	product, err := s.products.FindBySKU(sku)
	if err != nil {
		return kernel.Product{}, err
	}

	if !product.Active {
		return kernel.Product{}, ErrProductInactive
	}

	return kernel.Product{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		UnitPrice:        product.UnitPrice,
		ReturnWindowDays: product.ReturnWindowDays,
	}, nil
}

func (s Service) GetProduct(query kernel.GetProductQuery) (kernel.ProductDetails, error) {
	product, err := s.products.FindBySKU(query.SKU)
	if err != nil {
		return kernel.ProductDetails{}, err
	}

	return kernel.ProductDetails{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		Active:           product.Active,
		UnitPrice:        product.UnitPrice,
		ReturnWindowDays: product.ReturnWindowDays,
	}, nil
}

func (s Service) ListProducts(query kernel.ListProductsQuery) ([]kernel.ProductSummary, error) {
	productsList, err := s.products.List(query.Category, query.Active)
	if err != nil {
		return nil, err
	}

	results := make([]kernel.ProductSummary, 0, len(productsList))
	for _, product := range productsList {
		results = append(results, kernel.ProductSummary{
			SKU:              product.SKU,
			Name:             product.Name,
			Category:         product.Category,
			Active:           product.Active,
			UnitPrice:        product.UnitPrice,
			ReturnWindowDays: product.ReturnWindowDays,
		})
	}

	return results, nil
}
