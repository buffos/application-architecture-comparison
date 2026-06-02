package idempotency

type Result struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

type Store interface {
	Find(key string) (Result, bool, error)
	Save(key string, result Result) error
}
