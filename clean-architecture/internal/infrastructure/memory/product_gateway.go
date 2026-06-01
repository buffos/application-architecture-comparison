package memory

import (
	"sort"
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

func (g *ProductGateway) List(category string, availableOnly bool) ([]entities.Product, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	products := make([]entities.Product, 0, len(g.products))
	for _, product := range g.products {
		if category != "" && product.Category != category {
			continue
		}
		if availableOnly && !product.Available {
			continue
		}

		products = append(products, product)
	}

	sort.Slice(products, func(i int, j int) bool {
		return products[i].SKU < products[j].SKU
	})

	return products, nil
}
