package inventory

// Reserver is the public contract provided by this component to workflows
// that need to claim stock without accessing stock state directly.
type Reserver interface {
	Reserve(items []ReservationItem) error
}
