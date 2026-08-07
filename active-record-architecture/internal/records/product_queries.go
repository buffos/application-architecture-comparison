package records

import "sort"

// GetProduct is the named detail-query form of FindProduct.
func GetProduct(db *Database, sku string) (*Product, error) {
	return FindProduct(db, sku)
}

// ListProducts returns reconstructed product snapshots ordered by SKU. An
// empty category matches every category; activeOnly controls availability
// filtering.
func ListProducts(db *Database, category string, activeOnly bool) ([]*Product, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	skus := make([]string, 0, len(db.products))
	for sku, row := range db.products {
		if category != "" && row.Category != category {
			continue
		}
		if activeOnly && !row.Active {
			continue
		}
		skus = append(skus, sku)
	}
	sort.Strings(skus)

	products := make([]*Product, 0, len(skus))
	for _, sku := range skus {
		product, err := FindProduct(db, sku)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
