package scripts

import (
	"errors"

	"transaction-script-architecture/internal/data"
)

var (
	ErrQuoteIDRequired    = errors.New("quote id is required")
	ErrQuoteNotFound      = errors.New("quote not found")
	ErrQuoteNotEditable   = errors.New("quote is no longer editable")
	ErrProductSKURequired = errors.New("product sku is required")
	ErrQuantityInvalid    = errors.New("quantity must be positive")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductInactive    = errors.New("product is inactive")
)

// AddQuoteLine is a second transaction script. It coordinates the quote and
// product records directly, then saves the updated quote record.
func AddQuoteLine(store *data.Store, quoteID string, sku string, quantity int) (data.Quote, error) {
	if store == nil {
		return data.Quote{}, ErrStoreRequired
	}

	if quoteID == "" {
		return data.Quote{}, ErrQuoteIDRequired
	}

	if sku == "" {
		return data.Quote{}, ErrProductSKURequired
	}

	if quantity <= 0 {
		return data.Quote{}, ErrQuantityInvalid
	}

	quote, ok := store.Quotes[quoteID]
	if !ok {
		return data.Quote{}, ErrQuoteNotFound
	}

	if quote.Status != data.QuoteStatusDraft {
		return data.Quote{}, ErrQuoteNotEditable
	}

	product, ok := store.Products[sku]
	if !ok {
		return data.Quote{}, ErrProductNotFound
	}

	if !product.Active {
		return data.Quote{}, ErrProductInactive
	}

	quote.Lines = append(quote.Lines, data.QuoteLine{
		ProductCategory:     product.Category,
		SKU:                 product.SKU,
		ProductNameSnapshot: product.Name,
		Quantity:            quantity,
		UnitPrice:           product.UnitPrice,
		LineTotal:           product.UnitPrice * quantity,
	})

	store.Quotes[quote.ID] = quote

	return quote, nil
}
