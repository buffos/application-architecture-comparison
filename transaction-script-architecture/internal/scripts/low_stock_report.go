package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

type LowStockItem struct {
	SKU              string
	ProductName      string
	OnHand           int
	Reserved         int
	Available        int
	ReorderThreshold int
}

func GetLowStockItemsReport(store *data.Store) ([]LowStockItem, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	items := make([]LowStockItem, 0)
	for sku, stock := range store.Stocks {
		available := stock.OnHand - stock.Reserved
		if available > stock.ReorderThreshold {
			continue
		}

		item := LowStockItem{
			SKU:              sku,
			OnHand:           stock.OnHand,
			Reserved:         stock.Reserved,
			Available:        available,
			ReorderThreshold: stock.ReorderThreshold,
		}
		if product, ok := store.Products[sku]; ok {
			item.ProductName = product.Name
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})

	return items, nil
}
