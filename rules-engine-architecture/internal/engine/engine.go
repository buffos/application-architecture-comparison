package engine

import (
	"fmt"
	"sort"
)

// Engine owns registration and execution of the active Rule set.
type Engine struct {
	rules []Rule
}

func NewEngine() *Engine {
	return &Engine{rules: []Rule{}}
}

func (engine *Engine) Register(rule Rule) {
	engine.rules = append(engine.rules, rule)
}

func (engine *Engine) ExecuteAll(memory *WorkingMemory) error {
	orderedRules := append([]Rule(nil), engine.rules...)
	sort.SliceStable(orderedRules, func(left, right int) bool {
		return orderedRules[left].Priority() > orderedRules[right].Priority()
	})

	resolvedGroups := map[string]bool{}
	for _, rule := range orderedRules {
		group := rule.ConflictGroup()
		if group != "" && resolvedGroups[group] {
			continue
		}
		if !rule.Evaluate(memory) {
			continue
		}

		if err := rule.Execute(memory); err != nil {
			return fmt.Errorf("execute rule %q: %w", rule.Name(), err)
		}
		if group != "" {
			resolvedGroups[group] = true
		}
	}

	return nil
}
