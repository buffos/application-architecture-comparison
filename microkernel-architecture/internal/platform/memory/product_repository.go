package memory

import (
	"slices"
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

func (r *ProductRepository) List(category string, active *bool) ([]products.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]products.Product, 0)
	for _, product := range r.products {
		if category != "" && product.Category != category {
			continue
		}
		if active != nil && product.Active != *active {
			continue
		}
		results = append(results, product)
	}

	slices.SortFunc(results, func(a products.Product, b products.Product) int {
		if a.SKU < b.SKU {
			return -1
		}
		if a.SKU > b.SKU {
			return 1
		}
		return 0
	})

	return results, nil
}
