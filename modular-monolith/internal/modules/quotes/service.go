package quotes

import (
	"modular-monolith/internal/modules/approvals"
	"modular-monolith/internal/modules/products"
)

type CustomerDirectory interface {
	RequireActiveCustomer(id string) error
}

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type AddQuoteLineCommand struct {
	QuoteID    string
	ProductSKU string
	Quantity   int
}

type AddQuoteLineResult struct {
	QuoteID    string
	LineCount  int
	TotalItems int
	Status     string
}

type SubmitQuoteCommand struct {
	QuoteID string
}

type SubmitQuoteResult struct {
	QuoteID    string
	LineCount  int
	TotalItems int
	Status     string
}

type Service struct {
	quotes    Repository
	customers CustomerDirectory
	products  products.Catalog
	approvals approvals.Evaluator
}

func NewService(quotes Repository, customers CustomerDirectory, products products.Catalog, approvals approvals.Evaluator) Service {
	return Service{
		quotes:    quotes,
		customers: customers,
		products:  products,
		approvals: approvals,
	}
}

func (s Service) CreateDraftQuote(command CreateDraftQuoteCommand) (CreateDraftQuoteResult, error) {
	if err := s.customers.RequireActiveCustomer(command.CustomerID); err != nil {
		return CreateDraftQuoteResult{}, err
	}

	quote, err := NewDraftQuote(command.CustomerID)
	if err != nil {
		return CreateDraftQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return CreateDraftQuoteResult{}, err
	}

	return CreateDraftQuoteResult{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}, nil
}

func (s Service) AddQuoteLine(command AddQuoteLineCommand) (AddQuoteLineResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return AddQuoteLineResult{}, err
	}

	product, err := s.products.GetProductForQuote(command.ProductSKU)
	if err != nil {
		return AddQuoteLineResult{}, err
	}

	if err := quote.AddLine(ProductInput{
		SKU:       product.SKU,
		Name:      product.Name,
		Category:  product.Category,
		UnitPrice: product.UnitPrice,
	}, command.Quantity); err != nil {
		return AddQuoteLineResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return AddQuoteLineResult{}, err
	}

	return AddQuoteLineResult{
		QuoteID:    quote.ID,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
		Status:     quote.Status,
	}, nil
}

func (s Service) SubmitQuote(command SubmitQuoteCommand) (SubmitQuoteResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return SubmitQuoteResult{}, err
	}

	submission := approvals.QuoteSubmission{
		Lines: make([]approvals.QuoteSubmissionLine, 0, len(quote.Lines)),
	}
	for _, line := range quote.Lines {
		submission.Lines = append(submission.Lines, approvals.QuoteSubmissionLine{
			ProductCategory: line.ProductCategory,
		})
	}

	if err := quote.Submit(s.approvals.RequiresApproval(submission)); err != nil {
		return SubmitQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return SubmitQuoteResult{}, err
	}

	return SubmitQuoteResult{
		QuoteID:    quote.ID,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
		Status:     quote.Status,
	}, nil
}
