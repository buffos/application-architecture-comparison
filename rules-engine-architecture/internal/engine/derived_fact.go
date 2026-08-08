package engine

const ManagerApprovalRequiredFact = "manager-approval-required"

// DerivedFact is knowledge produced by one Rule for later Rules to consume.
type DerivedFact struct {
	Name  string
	Value string
}

func (memory *WorkingMemory) AddDerivedFact(fact DerivedFact) bool {
	for _, existing := range memory.DerivedFacts {
		if existing.Name == fact.Name && existing.Value == fact.Value {
			return false
		}
	}

	memory.DerivedFacts = append(memory.DerivedFacts, fact)
	return true
}

func (memory *WorkingMemory) HasDerivedFact(name string) bool {
	for _, fact := range memory.DerivedFacts {
		if fact.Name == name {
			return true
		}
	}
	return false
}
