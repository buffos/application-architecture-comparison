package data

// Payment is a passive record for a simulated payment attempt.
type Payment struct {
	ID              string
	OrderID         string
	Amount          int
	Status          string
	ReviewedBy      string
	DecisionComment string
}
