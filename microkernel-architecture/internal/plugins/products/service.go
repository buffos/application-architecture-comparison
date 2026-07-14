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
		SKU:       product.SKU,
		Name:      product.Name,
		Category:  product.Category,
		UnitPrice: product.UnitPrice,
	}, nil
}
