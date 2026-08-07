package quoting

import "errors"

var (
	ErrQuoteIDRequired        = errors.New("quote id is required")
	ErrCustomerIDRequired     = errors.New("customer id is required")
	ErrProductSKURequired     = errors.New("product sku is required")
	ErrQuantityMustBePositive = errors.New("quantity must be positive")
	ErrQuoteNotEditable       = errors.New("quote is no longer editable")
	ErrQuoteHasNoLines        = errors.New("quote must contain at least one line")
	ErrQuoteNotSubmittable    = errors.New("quote must be in draft status")
	ErrQuoteNotApprovable     = errors.New("quote must be pending approval")
	ErrQuoteNotRejectable     = errors.New("quote must be pending approval")
)

type QuoteID string
type CustomerID string
type ProductSKU string

type QuoteStatus string

const (
	QuoteStatusDraft           QuoteStatus = "Draft"
	QuoteStatusPendingApproval QuoteStatus = "PendingApproval"
	QuoteStatusApproved        QuoteStatus = "Approved"
	QuoteStatusRejected        QuoteStatus = "Rejected"
)

// QuoteLine is a validated part of a Quote aggregate. Its fields remain
// private so a caller cannot make a line invalid after construction.
type QuoteLine struct {
	sku       ProductSKU
	quantity  int
	unitPrice Money
}

func NewQuoteLine(sku ProductSKU, quantity int, unitPrice Money) (QuoteLine, error) {
	if sku == "" {
		return QuoteLine{}, ErrProductSKURequired
	}
	if quantity <= 0 {
		return QuoteLine{}, ErrQuantityMustBePositive
	}
	if unitPrice.Currency() == "" {
		return QuoteLine{}, ErrCurrencyRequired
	}

	return QuoteLine{sku: sku, quantity: quantity, unitPrice: unitPrice}, nil
}

func (line QuoteLine) ProductSKU() ProductSKU { return line.sku }
func (line QuoteLine) Quantity() int          { return line.quantity }
func (line QuoteLine) UnitPrice() Money       { return line.unitPrice }

func (line QuoteLine) Total() (Money, error) {
	return line.unitPrice.Multiply(line.quantity)
}

// Quote is the aggregate root for the Quoting domain. It owns the quote
// lifecycle and is deliberately unaware of persistence.
type Quote struct {
	id         QuoteID
	customerID CustomerID
	status     QuoteStatus
	lines      []QuoteLine
}

func NewQuote(id QuoteID, customerID CustomerID) (Quote, error) {
	if id == "" {
		return Quote{}, ErrQuoteIDRequired
	}
	if customerID == "" {
		return Quote{}, ErrCustomerIDRequired
	}

	return Quote{
		id:         id,
		customerID: customerID,
		status:     QuoteStatusDraft,
		lines:      []QuoteLine{},
	}, nil
}

func (quote Quote) ID() QuoteID            { return quote.id }
func (quote Quote) CustomerID() CustomerID { return quote.customerID }
func (quote Quote) Status() QuoteStatus    { return quote.status }
func (quote Quote) Lines() []QuoteLine     { return append([]QuoteLine(nil), quote.lines...) }
func (quote Quote) Currency() string {
	if len(quote.lines) == 0 {
		return ""
	}
	return quote.lines[0].unitPrice.Currency()
}

// AddLine is a business command on the aggregate, not a public slice append.
func (quote *Quote) AddLine(line QuoteLine) error {
	if quote.status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}
	if len(quote.lines) > 0 && quote.Currency() != line.unitPrice.Currency() {
		return ErrCurrencyMismatch
	}

	quote.lines = append(quote.lines, line)
	return nil
}

func (quote *Quote) SubmitForApproval() error {
	if quote.status != QuoteStatusDraft {
		return ErrQuoteNotSubmittable
	}
	if len(quote.lines) == 0 {
		return ErrQuoteHasNoLines
	}

	quote.status = QuoteStatusPendingApproval
	return nil
}

func (quote *Quote) Approve() error {
	if quote.status != QuoteStatusPendingApproval {
		return ErrQuoteNotApprovable
	}

	quote.status = QuoteStatusApproved
	return nil
}

func (quote *Quote) Reject() error {
	if quote.status != QuoteStatusPendingApproval {
		return ErrQuoteNotRejectable
	}

	quote.status = QuoteStatusRejected
	return nil
}

func (quote Quote) Total() (Money, error) {
	if len(quote.lines) == 0 {
		return Money{}, ErrQuoteHasNoLines
	}

	total, err := NewMoney(0, quote.Currency())
	if err != nil {
		return Money{}, err
	}
	for _, line := range quote.lines {
		lineTotal, err := line.Total()
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
