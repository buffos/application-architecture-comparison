package quoting

import "errors"

var ErrPricingTierInvalid = errors.New("pricing tier is invalid")

type CustomerPricingTier string

const (
	CustomerPricingTierStandard   CustomerPricingTier = "Standard"
	CustomerPricingTierPreferred  CustomerPricingTier = "Preferred"
	CustomerPricingTierEnterprise CustomerPricingTier = "Enterprise"
)

type ProductPricingInput struct {
	SKU         ProductSKU
	ProductName string
	Category    ProductCategory
	BasePrice   Money
}

type PricingPolicy interface {
	Price(ProductPricingInput, CustomerPricingTier) (Money, error)
}

type TierPricingPolicy struct{}

func (TierPricingPolicy) Price(product ProductPricingInput, tier CustomerPricingTier) (Money, error) {
	discountPercent, err := discountForTier(tier)
	if err != nil {
		return Money{}, err
	}

	adjustedCents := product.BasePrice.Cents() * int64(100-discountPercent) / 100
	return NewMoney(adjustedCents, product.BasePrice.Currency())
}

// QuotePricingService owns a stateless rule that spans Customer and Catalog
// facts but returns a value owned by the Quoting domain.
type QuotePricingService struct {
	policy PricingPolicy
}

func NewQuotePricingService() QuotePricingService {
	return NewQuotePricingServiceWithPolicy(TierPricingPolicy{})
}

func NewQuotePricingServiceWithPolicy(policy PricingPolicy) QuotePricingService {
	if policy == nil {
		policy = TierPricingPolicy{}
	}
	return QuotePricingService{policy: policy}
}

func (service QuotePricingService) PriceLine(product ProductPricingInput, tier CustomerPricingTier, quantity int) (QuoteLine, error) {
	price, err := service.policy.Price(product, tier)
	if err != nil {
		return QuoteLine{}, err
	}

	category := product.Category
	if category == "" {
		category = ProductCategoryStandard
	}
	if product.ProductName == "" {
		return NewQuoteLineWithCategory(product.SKU, category, quantity, price)
	}
	return NewQuoteLineFromProductSnapshotWithCategory(product.SKU, product.ProductName, category, quantity, price)
}

func discountForTier(tier CustomerPricingTier) (int64, error) {
	switch tier {
	case CustomerPricingTierStandard:
		return 0, nil
	case CustomerPricingTierPreferred:
		return 5, nil
	case CustomerPricingTierEnterprise:
		return 10, nil
	default:
		return 0, ErrPricingTierInvalid
	}
}
