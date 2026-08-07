package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

func GetQuote(store *data.Store, quoteID string) (data.Quote, error) {
	if store == nil {
		return data.Quote{}, ErrStoreRequired
	}
	if quoteID == "" {
		return data.Quote{}, ErrQuoteIDRequired
	}

	quote, ok := store.Quotes[quoteID]
	if !ok {
		return data.Quote{}, ErrQuoteNotFound
	}

	return cloneQuote(quote), nil
}

func ListQuotes(store *data.Store, status string, customerID string) ([]data.Quote, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	quotes := make([]data.Quote, 0, len(store.Quotes))
	for _, quote := range store.Quotes {
		if status != "" && quote.Status != status {
			continue
		}
		if customerID != "" && quote.CustomerID != customerID {
			continue
		}
		quotes = append(quotes, cloneQuote(quote))
	}

	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].ID < quotes[j].ID
	})

	return quotes, nil
}

func cloneQuote(quote data.Quote) data.Quote {
	quote.Lines = append([]data.QuoteLine(nil), quote.Lines...)
	return quote
}
