package application

import "onion-architecture/internal/domain"

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

type QuoteStore interface {
	FindByID(id string) (domain.Quote, error)
	Save(quote domain.Quote) error
}

type ProductLookup interface {
	FindBySKU(sku string) (domain.Product, error)
}

type AddQuoteLineService struct {
	quotes   QuoteStore
	products ProductLookup
}

func NewAddQuoteLineService(quotes QuoteStore, products ProductLookup) AddQuoteLineService {
	return AddQuoteLineService{
		quotes:   quotes,
		products: products,
	}
}

func (s AddQuoteLineService) Execute(command AddQuoteLineCommand) (AddQuoteLineResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return AddQuoteLineResult{}, err
	}

	product, err := s.products.FindBySKU(command.ProductSKU)
	if err != nil {
		return AddQuoteLineResult{}, err
	}

	if err := product.EnsureActive(); err != nil {
		return AddQuoteLineResult{}, err
	}

	if err := quote.AddLine(product, command.Quantity); err != nil {
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
