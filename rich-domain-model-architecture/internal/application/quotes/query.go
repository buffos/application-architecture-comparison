package quotes

import (
	"errors"
	"sort"

	"rich-domain-model-architecture/internal/domain/quoting"
)

var ErrQuoteNotFound = errors.New("quote not found")

type Reader interface {
	GetQuote(id string) (QuoteDetails, error)
	ListQuotes(status quoting.QuoteStatus) []QuoteSummary
}

type QuoteDetails struct {
	QuoteID    string
	CustomerID string
	Status     string
	TotalCents int64
	Currency   string
	LineCount  int
}

type QuoteSummary struct {
	QuoteID    string
	CustomerID string
	Status     string
	TotalCents int64
	LineCount  int
}

type InMemoryReader struct {
	quotes map[string]QuoteDetails
}

func NewInMemoryReader() *InMemoryReader {
	return &InMemoryReader{quotes: make(map[string]QuoteDetails)}
}

func (reader *InMemoryReader) Save(quote quoting.Quote) error {
	total, err := quote.Total()
	if err != nil {
		return err
	}
	details := QuoteDetails{
		QuoteID:    string(quote.ID()),
		CustomerID: string(quote.CustomerID()),
		Status:     string(quote.Status()),
		TotalCents: total.Cents(),
		Currency:   total.Currency(),
		LineCount:  len(quote.Lines()),
	}
	reader.quotes[details.QuoteID] = details
	return nil
}

func (reader *InMemoryReader) GetQuote(id string) (QuoteDetails, error) {
	details, ok := reader.quotes[id]
	if !ok {
		return QuoteDetails{}, ErrQuoteNotFound
	}
	return details, nil
}

func (reader *InMemoryReader) ListQuotes(status quoting.QuoteStatus) []QuoteSummary {
	result := make([]QuoteSummary, 0, len(reader.quotes))
	for _, details := range reader.quotes {
		if status != "" && details.Status != string(status) {
			continue
		}
		result = append(result, QuoteSummary{QuoteID: details.QuoteID, CustomerID: details.CustomerID, Status: details.Status, TotalCents: details.TotalCents, LineCount: details.LineCount})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].QuoteID < result[j].QuoteID })
	return result
}

var _ Reader = (*InMemoryReader)(nil)
