package quotes

import "microkernel-architecture/internal/kernel"

type Service struct {
	quotes    Repository
	customers kernel.CustomerDirectory
	products  kernel.ProductCatalog
}

func NewService(quotes Repository, customers kernel.CustomerDirectory, products kernel.ProductCatalog) Service {
	return Service{
		quotes:    quotes,
		customers: customers,
		products:  products,
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
