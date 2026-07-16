package quotes

import (
	"fmt"

	"component-based-architecture/internal/components/approvals"
	"component-based-architecture/internal/components/customers"
	"component-based-architecture/internal/components/products"
)

// Component owns quote behavior and its in-memory state for this lesson.
type Component struct {
	customers customers.CustomerDirectory
	approvals approvals.Evaluator
	products  products.Catalog
	quotes    map[string]Quote
	nextID    int
}

func NewComponent(customers customers.CustomerDirectory, products products.Catalog, approvals approvals.Evaluator) *Component {
	return &Component{
		customers: customers,
		approvals: approvals,
		products:  products,
		quotes:    make(map[string]Quote),
	}
}

type AddQuoteLineCommand struct {
	QuoteID    string
	ProductSKU string
	Quantity   int
}

type AddQuoteLineResult struct {
	QuoteID   string
	LineCount int
	Status    string
}

type SubmitQuoteCommand struct {
	QuoteID string
}

type SubmitQuoteResult struct {
	QuoteID   string
	LineCount int
	Status    string
}

type ApproveQuoteCommand struct {
	QuoteID string
}

type ApproveQuoteResult struct {
	QuoteID   string
	LineCount int
	Status    string
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

func (c *Component) AddQuoteLine(command AddQuoteLineCommand) (AddQuoteLineResult, error) {
	quote, ok := c.quotes[command.QuoteID]
	if !ok {
		return AddQuoteLineResult{}, ErrQuoteNotFound
	}

	product, err := c.products.GetProductForQuote(command.ProductSKU)
	if err != nil {
		return AddQuoteLineResult{}, err
	}
	if err := quote.AddLine(ProductInput{SKU: product.SKU, Name: product.Name, Category: product.Category, UnitPrice: product.UnitPrice}, command.Quantity); err != nil {
		return AddQuoteLineResult{}, err
	}
	c.quotes[quote.ID] = quote

	return AddQuoteLineResult{QuoteID: quote.ID, LineCount: len(quote.Lines), Status: quote.Status}, nil
}

func (c *Component) SubmitQuote(command SubmitQuoteCommand) (SubmitQuoteResult, error) {
	quote, ok := c.quotes[command.QuoteID]
	if !ok {
		return SubmitQuoteResult{}, ErrQuoteNotFound
	}
	submission := approvals.QuoteSubmission{Lines: make([]approvals.QuoteSubmissionLine, 0, len(quote.Lines))}
	for _, line := range quote.Lines {
		submission.Lines = append(submission.Lines, approvals.QuoteSubmissionLine{ProductCategory: line.ProductCategory})
	}
	if err := quote.Submit(c.approvals.RequiresApproval(submission)); err != nil {
		return SubmitQuoteResult{}, err
	}
	c.quotes[quote.ID] = quote

	return SubmitQuoteResult{QuoteID: quote.ID, LineCount: len(quote.Lines), Status: quote.Status}, nil
}

func (c *Component) ApproveQuote(command ApproveQuoteCommand) (ApproveQuoteResult, error) {
	quote, ok := c.quotes[command.QuoteID]
	if !ok {
		return ApproveQuoteResult{}, ErrQuoteNotFound
	}
	if err := quote.Approve(); err != nil {
		return ApproveQuoteResult{}, err
	}
	c.quotes[quote.ID] = quote

	return ApproveQuoteResult{QuoteID: quote.ID, LineCount: len(quote.Lines), Status: quote.Status}, nil
}
