package quotes

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrCustomerIDRequired = errors.New("customer id is required")
var ErrQuoteNotFound = errors.New("quote not found")
var ErrQuantityMustBePositive = errors.New("quantity must be positive")
var ErrQuoteNotEditable = errors.New("quote is not editable")
var ErrQuoteNotSubmittable = errors.New("quote is not submittable")
var ErrQuoteCannotBeSubmittedWithoutLines = errors.New("quote cannot be submitted without lines")
var ErrQuoteNotApprovable = errors.New("quote is not approvable")
var ErrQuoteNotConvertible = errors.New("quote is not convertible")

const QuoteStatusDraft = "Draft"
const QuoteStatusPendingApproval = "PendingApproval"
const QuoteStatusApproved = "Approved"

var quoteSequence uint64

type Quote struct {
	ID         string
	CustomerID string
	Status     string
	Lines      []QuoteLine
}

type QuoteLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

func NewDraftQuote(customerID string) (Quote, error) {
	if customerID == "" {
		return Quote{}, ErrCustomerIDRequired
	}

	id := atomic.AddUint64(&quoteSequence, 1)

	return Quote{
		ID:         fmt.Sprintf("quote-%03d", id),
		CustomerID: customerID,
		Status:     QuoteStatusDraft,
	}, nil
}

func (q *Quote) AddLine(product kernelProductInput, quantity int) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}

	if quantity <= 0 {
		return ErrQuantityMustBePositive
	}

	q.Lines = append(q.Lines, QuoteLine{
		ProductSKU:       product.SKU,
		ProductName:      product.Name,
		ProductCategory:  product.Category,
		Quantity:         quantity,
		UnitPrice:        product.UnitPrice,
		ReturnWindowDays: product.ReturnWindowDays,
	})

	return nil
}

func (q Quote) TotalQuantity() int {
	total := 0
	for _, line := range q.Lines {
		total += line.Quantity
	}

	return total
}

func (q *Quote) Submit(requiresApproval bool) error {
	if q.Status != QuoteStatusDraft {
		return ErrQuoteNotSubmittable
	}

	if len(q.Lines) == 0 {
		return ErrQuoteCannotBeSubmittedWithoutLines
	}

	if requiresApproval {
		q.Status = QuoteStatusPendingApproval
		return nil
	}

	q.Status = QuoteStatusApproved
	return nil
}

func (q *Quote) Approve() error {
	if q.Status != QuoteStatusPendingApproval {
		return ErrQuoteNotApprovable
	}

	q.Status = QuoteStatusApproved
	return nil
}

func (q Quote) EnsureConvertible() error {
	if q.Status != QuoteStatusApproved {
		return ErrQuoteNotConvertible
	}

	return nil
}

type kernelProductInput struct {
	SKU              string
	Name             string
	Category         string
	UnitPrice        int
	ReturnWindowDays int
}
