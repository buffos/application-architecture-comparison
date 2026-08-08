package engine

// CustomerFact describes the customer information relevant to future rules.
type CustomerFact struct {
	ID           string
	Name         string
	Tier         string
	InvoiceTerms bool
}

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentAccepted PaymentStatus = "accepted"
	PaymentFailed   PaymentStatus = "failed"
)

type PaymentFact struct {
	Status PaymentStatus
}

type ShipmentRequestFact struct {
	Requested bool
}

type OrderStatus string

const (
	OrderDraft     OrderStatus = "draft"
	OrderConfirmed OrderStatus = "confirmed"
	OrderShipped   OrderStatus = "shipped"
	OrderCancelled OrderStatus = "cancelled"
)

type OrderFact struct {
	ID     string
	Status OrderStatus
}

type CancellationRequestFact struct {
	Requested bool
}

type StockShortagePolicy string

const (
	StockShortageBackorder StockShortagePolicy = "backorder"
	StockShortageReject    StockShortagePolicy = "reject"
)

// ProductFact describes a product and its current inventory information.
type ProductFact struct {
	ID                string
	Name              string
	Category          string
	UnitPriceCents    int64
	AvailableQuantity int
	ShortagePolicy    StockShortagePolicy
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
	Customer     CustomerFact
	Quote        QuoteFact
	Products     []ProductFact
	Payment      PaymentFact
	Shipment     ShipmentRequestFact
	Order        OrderFact
	Cancellation CancellationRequestFact
	Findings     []Finding
	DerivedFacts []DerivedFact
	Trace        []RuleTrace
}

func NewWorkingMemory(customer CustomerFact, quote QuoteFact, products []ProductFact) *WorkingMemory {
	quoteCopy := quote
	quoteCopy.Lines = append([]QuoteLineFact(nil), quote.Lines...)

	return &WorkingMemory{
		Customer:     customer,
		Quote:        quoteCopy,
		Products:     append([]ProductFact(nil), products...),
		Findings:     []Finding{},
		DerivedFacts: []DerivedFact{},
		Trace:        []RuleTrace{},
	}
}

func (memory *WorkingMemory) AddFinding(finding Finding) {
	for _, existing := range memory.Findings {
		if existing.RuleName == finding.RuleName &&
			existing.Severity == finding.Severity &&
			existing.Message == finding.Message {
			return
		}
	}

	memory.Findings = append(memory.Findings, finding)
}

// ResetInferences starts a fresh inference session while preserving source Facts.
func (memory *WorkingMemory) ResetInferences() {
	memory.Findings = []Finding{}
	memory.DerivedFacts = []DerivedFact{}
	memory.Trace = []RuleTrace{}
}
