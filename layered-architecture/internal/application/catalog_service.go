package application

import "layered-architecture/internal/domain"

type ProductRepository interface {
	Save(product domain.Product) error
	FindBySKU(sku string) (domain.Product, error)
}

type CatalogService struct {
	repo ProductRepository
}

func NewCatalogService(repo ProductRepository) CatalogService {
	return CatalogService{repo: repo}
}

func (s CatalogService) CreateProduct(sku string, name string, category string, available bool) (domain.Product, error) {
	product, err := domain.NewProduct(sku, name, category, available)
	if err != nil {
		return domain.Product{}, err
	}

	if err := s.repo.Save(product); err != nil {
		return domain.Product{}, err
	}

	return product, nil
}
