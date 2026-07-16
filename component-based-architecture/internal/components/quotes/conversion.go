package quotes

// ApprovedQuoteSource is the public contract that this component provides to
// consumers that need a quote snapshot suitable for conversion.
type ApprovedQuoteSource interface {
	GetApprovedQuoteForOrder(quoteID string) (ApprovedQuote, error)
}

type ApprovedQuote struct {
	QuoteID    string
	CustomerID string
	Lines      []ApprovedQuoteLine
}

type ApprovedQuoteLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func (q Quote) EnsureConvertible() error {
	if q.Status != QuoteStatusApproved {
		return ErrQuoteNotConvertible
	}
	return nil
}

func (c *Component) GetApprovedQuoteForOrder(quoteID string) (ApprovedQuote, error) {
	quote, ok := c.quotes[quoteID]
	if !ok {
		return ApprovedQuote{}, ErrQuoteNotFound
	}
	if err := quote.EnsureConvertible(); err != nil {
		return ApprovedQuote{}, err
	}

	lines := make([]ApprovedQuoteLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, ApprovedQuoteLine{
			ProductSKU: line.ProductSKU, ProductName: line.ProductName, ProductCategory: line.ProductCategory,
			Quantity: line.Quantity, UnitPrice: line.UnitPrice, ReturnWindowDays: line.ReturnWindowDays,
		})
	}
	return ApprovedQuote{QuoteID: quote.ID, CustomerID: quote.CustomerID, Lines: lines}, nil
}

var _ ApprovedQuoteSource = (*Component)(nil)
