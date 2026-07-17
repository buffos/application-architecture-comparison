package quotes

// QuoteLookup is the public read contract that this component provides. It
// exposes a read model rather than the component's private storage.
type QuoteLookup interface {
	GetQuote(query GetQuoteQuery) (QuoteDetails, error)
	ListQuotes(query ListQuotesQuery) []QuoteSummary
}

type ListQuotesQuery struct {
	Status string
}

type QuoteSummary struct {
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type GetQuoteQuery struct {
	QuoteID string
}

type QuoteDetails struct {
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

func (c *Component) GetQuote(query GetQuoteQuery) (QuoteDetails, error) {
	quote, ok := c.quotes[query.QuoteID]
	if !ok {
		return QuoteDetails{}, ErrQuoteNotFound
	}

	return QuoteDetails{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
	}, nil
}

func (c *Component) ListQuotes(query ListQuotesQuery) []QuoteSummary {
	quotes := make([]QuoteSummary, 0, len(c.quotes))
	for _, quote := range c.quotes {
		if query.Status != "" && quote.Status != query.Status {
			continue
		}
		quotes = append(quotes, QuoteSummary{QuoteID: quote.ID, CustomerID: quote.CustomerID, Status: quote.Status, LineCount: len(quote.Lines)})
	}
	return quotes
}

var _ QuoteLookup = (*Component)(nil)
