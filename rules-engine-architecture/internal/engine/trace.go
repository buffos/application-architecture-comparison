package engine

// RuleTrace records one Rule's path through an Engine cycle.
type RuleTrace struct {
	RuleName      string
	Priority      int
	ConflictGroup string
	Enabled       bool
	Evaluated     bool
	Matched       bool
	Executed      bool
	SkippedReason string
}

func (memory *WorkingMemory) AddTrace(trace RuleTrace) {
	memory.Trace = append(memory.Trace, trace)
}
