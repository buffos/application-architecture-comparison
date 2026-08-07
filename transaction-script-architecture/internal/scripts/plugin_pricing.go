package scripts

import (
	"errors"
	"sort"
	"strconv"

	"transaction-script-architecture/internal/data"
)

const PluginDiscountPercentKey = "discountPercent"

var ErrPluginConfigurationInvalid = errors.New("pricing plugin configuration is invalid")

// PriceQuoteLine applies the enabled pricing plugin contributions in key
// order and returns the base unit price, total discount, and line total.
func PriceQuoteLine(store *data.Store, product data.Product, quantity int) (int, int, int, error) {
	if store == nil {
		return 0, 0, 0, ErrStoreRequired
	}
	if quantity <= 0 {
		return 0, 0, 0, ErrQuantityInvalid
	}

	keys := make([]string, 0, len(store.Plugins))
	for key, plugin := range store.Plugins {
		if plugin.Enabled && plugin.Type == PluginTypePricing {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	discountPercent := 0
	for _, key := range keys {
		plugin := store.Plugins[key]
		value := plugin.Config[PluginDiscountPercentKey]
		percent, err := strconv.Atoi(value)
		if err != nil || percent < 0 {
			return 0, 0, 0, ErrPluginConfigurationInvalid
		}
		discountPercent += percent
	}
	if discountPercent > 100 {
		discountPercent = 100
	}

	baseTotal := product.UnitPrice * quantity
	discountAmount := baseTotal * discountPercent / 100
	return product.UnitPrice, discountAmount, baseTotal - discountAmount, nil
}
