package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

func GetProduct(store *data.Store, sku string) (data.Product, error) {
	if store == nil {
		return data.Product{}, ErrStoreRequired
	}
	if sku == "" {
		return data.Product{}, ErrProductSKURequired
	}

	product, ok := store.Products[sku]
	if !ok {
		return data.Product{}, ErrProductNotFound
	}

	return product, nil
}

func ListProducts(store *data.Store, category string, activeOnly bool) ([]data.Product, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	products := make([]data.Product, 0, len(store.Products))
	for _, product := range store.Products {
		if category != "" && product.Category != category {
			continue
		}
		if activeOnly && !product.Active {
			continue
		}
		products = append(products, product)
	}

	sort.Slice(products, func(i, j int) bool {
		return products[i].SKU < products[j].SKU
	})

	return products, nil
}
