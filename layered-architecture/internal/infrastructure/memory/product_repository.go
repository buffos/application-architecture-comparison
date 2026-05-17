package memory

import (
	"sync"

	"layered-architecture/internal/domain"
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
