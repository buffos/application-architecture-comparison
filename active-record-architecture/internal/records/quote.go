package records

import "errors"

const (
	QuoteStatusDraft           = "Draft"
	QuoteStatusPendingApproval = "PendingApproval"
	QuoteStatusApproved        = "Approved"
	QuoteStatusRejected        = "Rejected"
)

var (
	ErrQuoteIDRequired         = errors.New("quote id is required")
	ErrQuoteCustomerIDRequired = errors.New("quote customer id is required")
	ErrQuoteStatusRequired     = errors.New("quote status is required")
	ErrQuoteNotFound           = errors.New("quote not found")
	ErrQuoteNotEditable        = errors.New("quote is no longer editable")
	ErrQuoteNotSubmittable     = errors.New("quote must be in draft status")
	ErrQuoteHasNoLines         = errors.New("quote must have at least one line")
	ErrReviewerRequired        = errors.New("reviewer is required")
	ErrQuoteNotApprovable      = errors.New("quote must be pending approval")
	ErrProductRequired         = errors.New("product is required")
	ErrQuantityInvalid         = errors.New("quantity must be positive")
)

const ApprovalReasonCustomBuild = "custom_build_requires_review"

// ApprovalDecision is a non-persistent result of evaluating a quote's stored
// line snapshots.
type ApprovalDecision struct {
	Required bool
	Reasons  []string
}

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

	ID               string
	CustomerID       string
	Status           string
	ConvertedOrderID string
	ReviewedBy       string
	DecisionComment  string
	Lines            []QuoteLine
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
		db:               db,
		ID:               row.ID,
		CustomerID:       row.CustomerID,
		Status:           row.Status,
		ConvertedOrderID: row.ConvertedOrderID,
		ReviewedBy:       row.ReviewedBy,
		DecisionComment:  row.DecisionComment,
		Lines:            cloneQuoteLines(row.Lines),
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

	decision := quote.EvaluateApproval()
	quote.Status = QuoteStatusApproved
	if decision.Required {
		quote.Status = QuoteStatusPendingApproval
	}
	return nil
}

// EvaluateApproval inspects the persisted quote snapshot without changing the
// Active Record. The returned reason codes are stable enough for later
// workflows and reports to reuse.
func (quote *Quote) EvaluateApproval() ApprovalDecision {
	decision := ApprovalDecision{}
	if quote == nil {
		return decision
	}
	for _, line := range quote.Lines {
		if line.ProductCategory != "CustomBuild" {
			continue
		}
		decision.Required = true
		decision.Reasons = append(decision.Reasons, ApprovalReasonCustomBuild)
		break
	}
	return decision
}

// Approve moves a pending quote to Approved and records the reviewer on the
// same Active Record row.
func (quote *Quote) Approve(reviewedBy string, decisionComment string) error {
	if quote == nil || quote.db == nil {
		return ErrDatabaseRequired
	}
	if reviewedBy == "" {
		return ErrReviewerRequired
	}
	if quote.Status != QuoteStatusPendingApproval {
		return ErrQuoteNotApprovable
	}

	quote.Status = QuoteStatusApproved
	quote.ReviewedBy = reviewedBy
	quote.DecisionComment = decisionComment
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
		ID:               quote.ID,
		CustomerID:       quote.CustomerID,
		Status:           quote.Status,
		ConvertedOrderID: quote.ConvertedOrderID,
		ReviewedBy:       quote.ReviewedBy,
		DecisionComment:  quote.DecisionComment,
		Lines:            cloneQuoteLines(quote.Lines),
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
