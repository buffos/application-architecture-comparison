package records

import "sort"

// LowStockItem is one low-stock read-model row.
type LowStockItem struct {
	SKU              string
	ProductName      string
	OnHand           int
	Reserved         int
	Available        int
	ReorderThreshold int
}

// GetLowStockItemsReport returns stock rows at or below their reorder
// threshold without changing inventory or catalog records.
func GetLowStockItemsReport(db *Database) ([]LowStockItem, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	items := make([]LowStockItem, 0)
	for sku, row := range db.stocks {
		available := row.OnHand - row.Reserved
		if available > row.ReorderThreshold {
			continue
		}

		item := LowStockItem{
			SKU:              sku,
			OnHand:           row.OnHand,
			Reserved:         row.Reserved,
			Available:        available,
			ReorderThreshold: row.ReorderThreshold,
		}
		if product, ok := db.products[sku]; ok {
			item.ProductName = product.Name
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})
	return items, nil
}
