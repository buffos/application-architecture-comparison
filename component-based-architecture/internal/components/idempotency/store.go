package idempotency

// Result is the completed return-review outcome that can be replayed safely.
type Result struct {
	ReturnRequestID string
	OrderID         string
	CustomerID      string
	Status          string
	LineCount       int
}

// Store is the public contract provided to workflows that need retry-safe
// command handling.
type Store interface {
	Find(key string) (Result, bool)
	Save(key string, result Result)
}
