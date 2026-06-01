package memory

import (
	"sync"

	"clean-architecture/internal/entities"
)

type ProductGateway struct {
	mu       sync.RWMutex
	products map[string]entities.Product
}

func NewProductGateway() *ProductGateway {
	return &ProductGateway{
		products: make(map[string]entities.Product),
	}
}

func (g *ProductGateway) Save(product entities.Product) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.products[product.SKU] = product
	return nil
}

func (g *ProductGateway) FindBySKU(sku string) (entities.Product, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	product, ok := g.products[sku]
	if !ok {
		return entities.Product{}, entities.ErrProductNotFound
	}

	return product, nil
}
