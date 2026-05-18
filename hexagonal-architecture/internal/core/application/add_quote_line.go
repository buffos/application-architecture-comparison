package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type AddQuoteLineUseCase struct {
	quotes   ports.QuoteRepository
	products ports.ProductLookup
	pricing  ports.PricingPolicy
}

func NewAddQuoteLineUseCase(quotes ports.QuoteRepository, products ports.ProductLookup, pricing ports.PricingPolicy) AddQuoteLineUseCase {
	return AddQuoteLineUseCase{
		quotes:   quotes,
		products: products,
		pricing:  pricing,
	}
}

func (uc AddQuoteLineUseCase) Execute(quoteID string, sku string, quantity int) (domain.Quote, error) {
	quote, err := uc.quotes.FindByID(quoteID)
	if err != nil {
		return domain.Quote{}, err
	}

	product, err := uc.products.FindBySKU(sku)
	if err != nil {
		return domain.Quote{}, err
	}

	if !product.Available {
		return domain.Quote{}, domain.ErrProductUnavailable
	}

	adjustedPrice, err := uc.pricing.Price(product, quantity)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.AddLine(product, quantity, adjustedPrice); err != nil {
		return domain.Quote{}, err
	}

	if err := uc.quotes.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
