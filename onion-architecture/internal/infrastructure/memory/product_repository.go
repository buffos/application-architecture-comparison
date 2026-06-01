package memory

import (
	"sync"

	"onion-architecture/internal/domain"
)

type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]domain.Product
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[string]domain.Product),
	}
}

func (r *ProductRepository) Save(product domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.SKU] = product
	return nil
}

func (r *ProductRepository) FindBySKU(sku string) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[sku]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}

	return product, nil
}

func (r *ProductRepository) List(category string, activeOnly bool) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Product, 0)
	for _, product := range r.products {
		if category != "" && product.Category != category {
			continue
		}

		if activeOnly && !product.Active {
			continue
		}

		result = append(result, product)
	}

	return result, nil
}
