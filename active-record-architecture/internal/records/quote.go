package records

import "errors"

const (
	QuoteStatusDraft           = "Draft"
	QuoteStatusPendingApproval = "PendingApproval"
	QuoteStatusApproved        = "Approved"
)

var (
	ErrQuoteIDRequired         = errors.New("quote id is required")
	ErrQuoteCustomerIDRequired = errors.New("quote customer id is required")
	ErrQuoteStatusRequired     = errors.New("quote status is required")
	ErrQuoteNotFound           = errors.New("quote not found")
	ErrQuoteNotEditable        = errors.New("quote is no longer editable")
	ErrQuoteNotSubmittable     = errors.New("quote must be in draft status")
	ErrQuoteHasNoLines         = errors.New("quote must have at least one line")
	ErrProductRequired         = errors.New("product is required")
	ErrQuantityInvalid         = errors.New("quantity must be positive")
)

// QuoteLine is a product snapshot embedded in a Quote Active Record.
type QuoteLine struct {
	ProductCategory     string
	SKU                 string
	ProductNameSnapshot string
	Quantity            int
	UnitPrice           int
	DiscountAmount      int
	ReturnWindowDays    int
	LineTotal           int
}

// Quote is an Active Record. It contains the quote fields and the database
// connection used by Save and FindQuote.
type Quote struct {
	db *Database

	ID         string
	CustomerID string
	Status     string
	Lines      []QuoteLine
}

// NewDraftQuote creates an unsaved draft quote for an existing active
// customer. The caller explicitly persists the returned Active Record with
// Quote.Save.
func NewDraftQuote(db *Database, customerID string) (*Quote, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if customerID == "" {
		return nil, ErrCustomerIDRequired
	}

	customer, err := FindCustomer(db, customerID)
	if err != nil {
		return nil, err
	}
	if !customer.Active {
		return nil, ErrCustomerInactive
	}

	return &Quote{
		db:         db,
		ID:         db.nextQuoteID(),
		CustomerID: customer.ID,
		Status:     QuoteStatusDraft,
		Lines:      []QuoteLine{},
	}, nil
}

// FindQuote loads a Quote Active Record from the quote table.
func FindQuote(db *Database, id string) (*Quote, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if id == "" {
		return nil, ErrQuoteIDRequired
	}

	row, ok := db.quotes[id]
	if !ok {
		return nil, ErrQuoteNotFound
	}

	return &Quote{
		db:         db,
		ID:         row.ID,
		CustomerID: row.CustomerID,
		Status:     row.Status,
		Lines:      cloneQuoteLines(row.Lines),
	}, nil
}

// AddLine applies the first quote-editing behavior directly on the Active
// Record. The product is copied into the line so later catalog changes do not
// rewrite the commercial snapshot.
func (quote *Quote) AddLine(product *Product, quantity int) error {
	if quote == nil || quote.db == nil {
		return ErrDatabaseRequired
	}
	if quote.Status != QuoteStatusDraft {
		return ErrQuoteNotEditable
	}
	if product == nil {
		return ErrProductRequired
	}
	if quantity <= 0 {
		return ErrQuantityInvalid
	}
	if !product.Active {
		return ErrProductInactive
	}

	quote.Lines = append(quote.Lines, QuoteLine{
		ProductCategory:     product.Category,
		SKU:                 product.SKU,
		ProductNameSnapshot: product.Name,
		Quantity:            quantity,
		UnitPrice:           product.UnitPrice,
		ReturnWindowDays:    product.ReturnWindowDays,
		LineTotal:           product.UnitPrice * quantity,
	})
	return nil
}

// SubmitForApproval applies the first quote lifecycle decision directly on
// the Active Record. Custom-build lines require a review; other quotes can be
// approved automatically.
func (quote *Quote) SubmitForApproval() error {
	if quote == nil || quote.db == nil {
		return ErrDatabaseRequired
	}
	if quote.Status != QuoteStatusDraft {
		return ErrQuoteNotSubmittable
	}
	if len(quote.Lines) == 0 {
		return ErrQuoteHasNoLines
	}

	quote.Status = QuoteStatusApproved
	for _, line := range quote.Lines {
		if line.ProductCategory == "CustomBuild" {
			quote.Status = QuoteStatusPendingApproval
			break
		}
	}
	return nil
}

// Save writes the current Quote Active Record to its table.
func (quote *Quote) Save() error {
	if quote == nil || quote.db == nil {
		return ErrDatabaseRequired
	}
	if quote.ID == "" {
		return ErrQuoteIDRequired
	}
	if quote.CustomerID == "" {
		return ErrQuoteCustomerIDRequired
	}
	if quote.Status == "" {
		return ErrQuoteStatusRequired
	}
	if _, ok := quote.db.customers[quote.CustomerID]; !ok {
		return ErrCustomerNotFound
	}

	quote.db.quotes[quote.ID] = quoteRow{
		ID:         quote.ID,
		CustomerID: quote.CustomerID,
		Status:     quote.Status,
		Lines:      cloneQuoteLines(quote.Lines),
	}
	return nil
}

func cloneQuoteLines(lines []QuoteLine) []QuoteLine {
	if lines == nil {
		return nil
	}
	clone := make([]QuoteLine, len(lines))
	copy(clone, lines)
	return clone
}
