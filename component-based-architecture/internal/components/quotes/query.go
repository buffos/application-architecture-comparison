package quotes

// QuoteLookup is the public read contract that this component provides. It
// exposes a read model rather than the component's private storage.
type QuoteLookup interface {
	GetQuote(query GetQuoteQuery) (QuoteDetails, error)
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

var _ QuoteLookup = (*Component)(nil)
