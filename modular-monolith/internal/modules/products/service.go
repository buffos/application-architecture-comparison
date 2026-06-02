package products

type Service struct {
	products Repository
}

func NewService(products Repository) Service {
	return Service{
		products: products,
	}
}

func (s Service) GetProductForQuote(sku string) (ProductForQuote, error) {
	product, err := s.products.FindBySKU(sku)
	if err != nil {
		return ProductForQuote{}, err
	}

	if !product.Active {
		return ProductForQuote{}, ErrProductInactive
	}

	return ProductForQuote{
		SKU:              product.SKU,
		Name:             product.Name,
		Category:         product.Category,
		UnitPrice:        product.UnitPrice,
		ReturnWindowDays: product.ReturnWindowDays,
	}, nil
}
