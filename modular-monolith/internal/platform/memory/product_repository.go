package memory

import (
	"sync"

	"modular-monolith/internal/modules/products"
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

func (r *ProductRepository) List(category string, activeOnly bool) ([]products.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]products.Product, 0, len(r.products))
	for _, product := range r.products {
		if category != "" && product.Category != category {
			continue
		}
		if activeOnly && !product.Active {
			continue
		}
		list = append(list, product)
	}

	return list, nil
}
