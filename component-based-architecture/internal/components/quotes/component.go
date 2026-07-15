package quotes

import (
	"fmt"

	"component-based-architecture/internal/components/customers"
)

// Component owns quote behavior and its in-memory state for this lesson.
type Component struct {
	customers customers.CustomerDirectory
	quotes    map[string]Quote
	nextID    int
}

func NewComponent(customers customers.CustomerDirectory) *Component {
	return &Component{
		customers: customers,
		quotes:    make(map[string]Quote),
	}
}

func (c *Component) CreateDraftQuote(command CreateDraftQuoteCommand) (CreateDraftQuoteResult, error) {
	if command.CustomerID == "" {
		return CreateDraftQuoteResult{}, ErrCustomerIDRequired
	}

	if err := c.customers.RequireActiveCustomer(command.CustomerID); err != nil {
		return CreateDraftQuoteResult{}, err
	}

	c.nextID++
	quote := Quote{
		ID:         fmt.Sprintf("quote-%03d", c.nextID),
		CustomerID: command.CustomerID,
		Status:     QuoteStatusDraft,
	}
	c.quotes[quote.ID] = quote

	return CreateDraftQuoteResult{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}, nil
}
