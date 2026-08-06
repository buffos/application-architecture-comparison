package quoting

import "errors"

var (
	ErrQuoteIDRequired        = errors.New("quote id is required")
	ErrCustomerIDRequired     = errors.New("customer id is required")
	ErrProductSKURequired     = errors.New("product sku is required")
	ErrProductCategoryInvalid = errors.New("product category is invalid")
	ErrQuantityMustBePositive = errors.New("quantity must be positive")
	ErrQuoteNotEditable       = errors.New("quote is not editable")
	ErrQuoteHasNoLines        = errors.New("quote must contain at least one line")
)

type QuoteID string
type CustomerID string
type ProductSKU string
type ProductCategory string

const (
	ProductCategoryStandard    ProductCategory = "Standard"
	ProductCategoryCustomBuild ProductCategory = "CustomBuild"
	ProductCategoryClearance   ProductCategory = "Clearance"
)

type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "Draft"
	QuoteStatusSubmitted QuoteStatus = "Submitted"
)

type QuoteLine struct {
	sku       ProductSKU
	category  ProductCategory
	quantity  int
	unitPrice Money
}

func NewQuoteLine(sku ProductSKU, quantity int, unitPrice Money) (QuoteLine, error) {
	return NewQuoteLineWithCategory(sku, ProductCategoryStandard, quantity, unitPrice)
}

func NewQuoteLineWithCategory(sku ProductSKU, category ProductCategory, quantity int, unitPrice Money) (QuoteLine, error) {
	if sku == "" {
		return QuoteLine{}, ErrProductSKURequired
	}
	if !validProductCategory(category) {
		return QuoteLine{}, ErrProductCategoryInvalid
	}
	if quantity <= 0 {
		return QuoteLine{}, ErrQuantityMustBePositive
	}
	if unitPrice.Currency() == "" {
		return QuoteLine{}, ErrCurrencyRequired
	}
	return QuoteLine{sku: sku, category: category, quantity: quantity, unitPrice: unitPrice}, nil
}

func (line QuoteLine) ProductSKU() ProductSKU           { return line.sku }
func (line QuoteLine) ProductCategory() ProductCategory { return line.category }
func (line QuoteLine) Quantity() int                    { return line.quantity }
func (line QuoteLine) UnitPrice() Money                 { return line.unitPrice }

// Quote is the aggregate root for the Quoting bounded context.
type Quote struct {
	id         QuoteID
	customerID CustomerID
	status     QuoteStatus
	currency   string
	lines      []QuoteLine
}

func NewQuote(id QuoteID, customerID CustomerID) (Quote, error) {
	if id == "" {
		return Quote{}, ErrQuoteIDRequired
	}
	if customerID == "" {
		return Quote{}, ErrCustomerIDRequired
	}
	return Quote{id: id, customerID: customerID, status: QuoteStatusDraft}, nil
}

func (q Quote) ID() QuoteID            { return q.id }
func (q Quote) CustomerID() CustomerID { return q.customerID }
func (q Quote) Status() QuoteStatus    { return q.status }
func (q Quote) Lines() []QuoteLine     { return append([]QuoteLine(nil), q.lines...) }
func (q Quote) Currency() string       { return q.currency }

func (q *Quote) AddLine(line QuoteLine) error {
	if q.status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}
	if q.currency != "" && q.currency != line.unitPrice.Currency() {
		return ErrCurrencyMismatch
	}
	if q.currency == "" {
		q.currency = line.unitPrice.Currency()
	}
	q.lines = append(q.lines, line)
	return nil
}

func (q *Quote) Submit() error {
	if q.status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}
	if len(q.lines) == 0 {
		return ErrQuoteHasNoLines
	}
	q.status = QuoteStatusSubmitted
	return nil
}

func (q Quote) Total() (Money, error) {
	if len(q.lines) == 0 {
		return Money{}, ErrQuoteHasNoLines
	}
	total, err := NewMoney(0, q.currency)
	if err != nil {
		return Money{}, err
	}
	for _, line := range q.lines {
		lineTotal, err := line.unitPrice.Multiply(line.quantity)
		if err != nil {
			return Money{}, err
		}
		total, err = total.Add(lineTotal)
		if err != nil {
			return Money{}, err
		}
	}
	return total, nil
}

func validProductCategory(category ProductCategory) bool {
	switch category {
	case ProductCategoryStandard, ProductCategoryCustomBuild, ProductCategoryClearance:
		return true
	default:
		return false
	}
}
