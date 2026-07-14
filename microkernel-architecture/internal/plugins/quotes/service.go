package quotes

import "microkernel-architecture/internal/kernel"

type Service struct {
	quotes    Repository
	customers kernel.CustomerDirectory
	products  kernel.ProductCatalog
	approvals kernel.ApprovalPolicy
}

func NewService(quotes Repository, customers kernel.CustomerDirectory, products kernel.ProductCatalog, approvals kernel.ApprovalPolicy) Service {
	return Service{
		quotes:    quotes,
		customers: customers,
		products:  products,
		approvals: approvals,
	}
}

func (s Service) CreateDraftQuote(command kernel.CreateDraftQuoteCommand) (kernel.CreateDraftQuoteResult, error) {
	if err := s.customers.RequireActiveCustomer(command.CustomerID); err != nil {
		return kernel.CreateDraftQuoteResult{}, err
	}

	quote, err := NewDraftQuote(command.CustomerID)
	if err != nil {
		return kernel.CreateDraftQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return kernel.CreateDraftQuoteResult{}, err
	}

	return kernel.CreateDraftQuoteResult{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
	}, nil
}

func (s Service) GetQuote(query kernel.GetQuoteQuery) (kernel.QuoteDetails, error) {
	quote, err := s.quotes.FindByID(query.QuoteID)
	if err != nil {
		return kernel.QuoteDetails{}, err
	}

	return kernel.QuoteDetails{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
	}, nil
}

func (s Service) AddQuoteLine(command kernel.AddQuoteLineCommand) (kernel.AddQuoteLineResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return kernel.AddQuoteLineResult{}, err
	}

	product, err := s.products.GetProductForQuote(command.ProductSKU)
	if err != nil {
		return kernel.AddQuoteLineResult{}, err
	}

	if err := quote.AddLine(kernelProductInput{
		SKU:       product.SKU,
		Name:      product.Name,
		Category:  product.Category,
		UnitPrice: product.UnitPrice,
	}, command.Quantity); err != nil {
		return kernel.AddQuoteLineResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return kernel.AddQuoteLineResult{}, err
	}

	return kernel.AddQuoteLineResult{
		QuoteID:    quote.ID,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
		Status:     quote.Status,
	}, nil
}

func (s Service) SubmitQuote(command kernel.SubmitQuoteCommand) (kernel.SubmitQuoteResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return kernel.SubmitQuoteResult{}, err
	}

	submission := kernel.QuoteSubmission{
		Lines: make([]kernel.QuoteSubmissionLine, 0, len(quote.Lines)),
	}
	for _, line := range quote.Lines {
		submission.Lines = append(submission.Lines, kernel.QuoteSubmissionLine{
			ProductCategory: line.ProductCategory,
		})
	}

	if err := quote.Submit(s.approvals.RequiresApproval(submission)); err != nil {
		return kernel.SubmitQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return kernel.SubmitQuoteResult{}, err
	}

	return kernel.SubmitQuoteResult{
		QuoteID:    quote.ID,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
		Status:     quote.Status,
	}, nil
}

func (s Service) ApproveQuote(command kernel.ApproveQuoteCommand) (kernel.ApproveQuoteResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return kernel.ApproveQuoteResult{}, err
	}

	if err := quote.Approve(); err != nil {
		return kernel.ApproveQuoteResult{}, err
	}

	if err := s.quotes.Save(quote); err != nil {
		return kernel.ApproveQuoteResult{}, err
	}

	return kernel.ApproveQuoteResult{
		QuoteID:    quote.ID,
		LineCount:  len(quote.Lines),
		TotalItems: quote.TotalQuantity(),
		Status:     quote.Status,
	}, nil
}

func (s Service) GetApprovedQuoteForOrder(quoteID string) (kernel.ApprovedQuote, error) {
	quote, err := s.quotes.FindByID(quoteID)
	if err != nil {
		return kernel.ApprovedQuote{}, err
	}

	if err := quote.EnsureConvertible(); err != nil {
		return kernel.ApprovedQuote{}, err
	}

	lines := make([]kernel.ApprovedQuoteLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, kernel.ApprovedQuoteLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
		})
	}

	return kernel.ApprovedQuote{
		QuoteID:    quote.ID,
		CustomerID: quote.CustomerID,
		Lines:      lines,
	}, nil
}
