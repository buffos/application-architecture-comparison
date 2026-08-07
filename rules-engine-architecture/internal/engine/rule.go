package engine

// Rule is the contract implemented by each independent business policy.
type Rule interface {
	Name() string
	Priority() int
	ConflictGroup() string
	Evaluate(memory *WorkingMemory) bool
	Execute(memory *WorkingMemory) error
}
