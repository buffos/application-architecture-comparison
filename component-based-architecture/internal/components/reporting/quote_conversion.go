package reporting

import (
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/quotes"
)

type QuoteConversionReport struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

type Component struct {
	quotes quotes.QuoteLookup
	orders orders.Reader
}

func NewComponent(quotes quotes.QuoteLookup, orders orders.Reader) *Component {
	return &Component{quotes: quotes, orders: orders}
}

func (c *Component) QuoteConversionReport() QuoteConversionReport {
	allQuotes := c.quotes.ListQuotes(quotes.ListQuotesQuery{})
	approvedQuotes := c.quotes.ListQuotes(quotes.ListQuotesQuery{Status: quotes.QuoteStatusApproved})
	converted := c.orders.ListOrders(orders.ListOrdersQuery{})
	convertedByQuote := make(map[string]struct{}, len(converted))
	for _, order := range converted {
		if order.QuoteID != "" {
			convertedByQuote[order.QuoteID] = struct{}{}
		}
	}
	result := QuoteConversionReport{TotalQuotes: len(allQuotes), ApprovedQuotes: len(approvedQuotes), ConvertedQuotes: len(convertedByQuote)}
	if result.TotalQuotes > 0 {
		result.ConversionRate = float64(result.ConvertedQuotes) / float64(result.TotalQuotes)
	}
	return result
}
