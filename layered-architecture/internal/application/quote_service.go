package application

import "layered-architecture/internal/domain"

type QuoteRepository interface {
	Save(quote domain.Quote) error
	FindByID(id string) (domain.Quote, error)
	List() ([]domain.Quote, error)
}

type QuoteService struct {
	repo           QuoteRepository
	customerRepo   CustomerRepository
	productRepo    ProductRepository
	pluginRegistry PricingPluginRegistry
}

func NewQuoteService(repo QuoteRepository, customerRepo CustomerRepository, productRepo ProductRepository, pluginRegistry PricingPluginRegistry) QuoteService {
	if pluginRegistry == nil {
		pluginRegistry = NoopPricingPluginRegistry{}
	}

	return QuoteService{
		repo:           repo,
		customerRepo:   customerRepo,
		productRepo:    productRepo,
		pluginRegistry: pluginRegistry,
	}
}

func (s QuoteService) CreateDraftQuote(customerID string) (domain.Quote, error) {
	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return domain.Quote{}, err
	}

	if !customer.Active {
		return domain.Quote{}, domain.ErrCustomerInactive
	}

	quote, err := domain.NewDraftQuote(customerID)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}

func (s QuoteService) GetQuote(id string) (domain.Quote, error) {
	return s.repo.FindByID(id)
}

func (s QuoteService) AddQuoteLine(id string, sku string, quantity int) (domain.Quote, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	product, err := s.productRepo.FindBySKU(sku)
	if err != nil {
		return domain.Quote{}, err
	}

	if !product.Available {
		return domain.Quote{}, domain.ErrProductUnavailable
	}

	adjustedPrice := product.BasePrice
	adjustments := make([]string, 0)
	for _, plugin := range s.pluginRegistry.EnabledPricingPlugins() {
		adjustment, ok := plugin.Adjust(PricingPluginInput{
			SKU:       product.SKU,
			Category:  product.Category,
			Quantity:  quantity,
			BasePrice: adjustedPrice,
		})
		if !ok {
			continue
		}

		adjustedPrice = adjustment.AdjustedPrice
		adjustments = append(adjustments, adjustment.Label)
	}

	if err := quote.AddLine(product.SKU, product.Category, product.Name, quantity, product.BasePrice, adjustedPrice, adjustments); err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}

func (s QuoteService) SubmitQuote(id string) (domain.Quote, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.Submit(); err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}

func (s QuoteService) ApproveQuote(id string) (domain.Quote, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.Approve(); err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}

func (s QuoteService) RejectQuote(id string) (domain.Quote, error) {
	quote, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Quote{}, err
	}

	if err := quote.Reject(); err != nil {
		return domain.Quote{}, err
	}

	if err := s.repo.Save(quote); err != nil {
		return domain.Quote{}, err
	}

	return quote, nil
}
