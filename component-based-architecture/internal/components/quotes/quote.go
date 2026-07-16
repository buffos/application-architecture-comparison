package quotes

import "errors"

var (
	ErrCustomerIDRequired                 = errors.New("customer id is required")
	ErrQuoteNotFound                      = errors.New("quote not found")
	ErrQuantityMustBePositive             = errors.New("quantity must be positive")
	ErrQuoteNotEditable                   = errors.New("quote is not editable")
	ErrQuoteNotSubmittable                = errors.New("quote is not submittable")
	ErrQuoteCannotBeSubmittedWithoutLines = errors.New("quote cannot be submitted without lines")
)

const (
	QuoteStatusDraft    = "Draft"
	QuoteStatusApproved = "Approved"
)

type Quote struct {
	ID         string
	CustomerID string
	Status     string
	Lines      []QuoteLine
}

type QuoteLine struct {
	ProductSKU  string
	ProductName string
	Quantity    int
	UnitPrice   int
}

type ProductInput struct {
	SKU       string
	Name      string
	UnitPrice int
}

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}

func (q *Quote) AddLine(product ProductInput, quantity int) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}
	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}
	q.Lines = append(q.Lines, QuoteLine{
		ProductSKU: product.SKU, ProductName: product.Name, Quantity: quantity, UnitPrice: product.UnitPrice,
	})
	return nil
}

func (q *Quote) Submit() error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteNotSubmittable
	}
	if len(q.Lines) == 0 {
		return ErrQuoteCannotBeSubmittedWithoutLines
	}
	q.Status = QuoteStatusApproved
	return nil
}
