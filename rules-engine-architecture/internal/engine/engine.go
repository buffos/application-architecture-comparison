package engine

import "fmt"

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
	for _, rule := range engine.rules {
		if !rule.Evaluate(memory) {
			continue
		}

		if err := rule.Execute(memory); err != nil {
			return fmt.Errorf("execute rule %q: %w", rule.Name(), err)
		}
	}

	return nil
}
