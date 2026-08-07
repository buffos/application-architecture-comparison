package engine

// CustomerFact describes the customer information relevant to future rules.
type CustomerFact struct {
	ID           string
	Name         string
	Tier         string
	InvoiceTerms bool
}

// ProductFact describes a product and its current inventory information.
type ProductFact struct {
	ID                string
	Name              string
	Category          string
	UnitPriceCents    int64
	AvailableQuantity int
}

// QuoteLineFact is a data-only line inside a QuoteFact.
type QuoteLineFact struct {
	ProductID      string
	Quantity       int
	UnitPriceCents int64
}

// QuoteFact describes the current quote state without owning business behavior.
type QuoteFact struct {
	ID              string
	CustomerID      string
	Lines           []QuoteLineFact
	DiscountPercent int
	Status          string
}

// Finding is a future output of a Rule evaluation.
type Finding struct {
	RuleName string
	Severity string
	Message  string
}

// WorkingMemory is the shared fact container used by the future Rule Engine.
type WorkingMemory struct {
	Customer CustomerFact
	Quote    QuoteFact
	Products []ProductFact
	Findings []Finding
}

func NewWorkingMemory(customer CustomerFact, quote QuoteFact, products []ProductFact) *WorkingMemory {
	quoteCopy := quote
	quoteCopy.Lines = append([]QuoteLineFact(nil), quote.Lines...)

	return &WorkingMemory{
		Customer: customer,
		Quote:    quoteCopy,
		Products: append([]ProductFact(nil), products...),
		Findings: []Finding{},
	}
}
