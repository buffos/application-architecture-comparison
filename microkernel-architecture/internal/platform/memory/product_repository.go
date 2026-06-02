package memory

import (
	"sync"

	"microkernel-architecture/internal/plugins/products"
)

type ProductRepository struct {
	mu       sync.RWMutex
	products map[string]products.Product
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[string]products.Product),
	}
}

func (r *ProductRepository) Save(product products.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.SKU] = product
	return nil
}

func (r *ProductRepository) FindBySKU(sku string) (products.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[sku]
	if !ok {
		return products.Product{}, products.ErrProductNotFound
	}

	return product, nil
}
