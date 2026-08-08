package engine

import (
	"fmt"
	"sort"
)

// Engine owns registration and execution of the active Rule set.
type Engine struct {
	rules []registeredRule
}

type registeredRule struct {
	rule    Rule
	enabled bool
}

func NewEngine() *Engine {
	return &Engine{rules: []registeredRule{}}
}

func (engine *Engine) Register(rule Rule) {
	engine.rules = append(engine.rules, registeredRule{
		rule:    rule,
		enabled: true,
	})
}

func (engine *Engine) SetRuleEnabled(ruleName string, enabled bool) bool {
	for index := range engine.rules {
		if engine.rules[index].rule.Name() == ruleName {
			engine.rules[index].enabled = enabled
			return true
		}
	}
	return false
}

func (engine *Engine) ExecuteAll(memory *WorkingMemory) error {
	orderedRules := append([]registeredRule(nil), engine.rules...)
	sort.SliceStable(orderedRules, func(left, right int) bool {
		return orderedRules[left].rule.Priority() > orderedRules[right].rule.Priority()
	})

	resolvedGroups := map[string]bool{}
	for _, registered := range orderedRules {
		if !registered.enabled {
			continue
		}

		rule := registered.rule
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
