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
	SKU       ProductSKU
	Category  ProductCategory
	BasePrice Money
}

// QuotePricingService is a stateless domain service for a rule that spans
// customer and catalog facts but produces a Quoting value.
type QuotePricingService struct{}

func NewQuotePricingService() QuotePricingService { return QuotePricingService{} }

func (QuotePricingService) PriceLine(product ProductPricingInput, tier CustomerPricingTier, quantity int) (QuoteLine, error) {
	discountPercent, err := discountForTier(tier)
	if err != nil {
		return QuoteLine{}, err
	}
	adjustedCents := product.BasePrice.Cents() * int64(100-discountPercent) / 100
	price, err := NewMoney(adjustedCents, product.BasePrice.Currency())
	if err != nil {
		return QuoteLine{}, err
	}
	category := product.Category
	if category == "" {
		category = ProductCategoryStandard
	}
	return NewQuoteLineWithCategory(product.SKU, category, quantity, price)
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
