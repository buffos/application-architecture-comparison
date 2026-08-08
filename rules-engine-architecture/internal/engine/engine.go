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
		rule := registered.rule
		group := rule.ConflictGroup()
		trace := RuleTrace{
			RuleName:      rule.Name(),
			Priority:      rule.Priority(),
			ConflictGroup: group,
			Enabled:       registered.enabled,
		}

		if !registered.enabled {
			trace.SkippedReason = "disabled by configuration"
			memory.AddTrace(trace)
			continue
		}

		if group != "" && resolvedGroups[group] {
			trace.SkippedReason = "conflict group already resolved"
			memory.AddTrace(trace)
			continue
		}

		trace.Evaluated = true
		if !rule.Evaluate(memory) {
			trace.SkippedReason = "condition did not match"
			memory.AddTrace(trace)
			continue
		}

		trace.Matched = true

		if err := rule.Execute(memory); err != nil {
			trace.SkippedReason = "execution failed"
			memory.AddTrace(trace)
			return fmt.Errorf("execute rule %q: %w", rule.Name(), err)
		}

		trace.Executed = true
		memory.AddTrace(trace)
		if group != "" {
			resolvedGroups[group] = true
		}
	}

	return nil
}
