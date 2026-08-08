package main

import (
	"flag"
	"fmt"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/rules"
)

func formatCents(amountCents int64) string {
	return fmt.Sprintf("%d.%02d", amountCents/100, amountCents%100)
}

func main() {
	disableCustomBuild := flag.Bool(
		"disable-custom-build",
		false,
		"disable the CustomBuild approval Rule",
	)
	simulateQuoteEdit := flag.Bool(
		"simulate-quote-edit",
		false,
		"edit the quote to a Standard product and recompute the decision",
	)
	flag.Parse()

	customer := engine.CustomerFact{
		ID:           "CUST-001",
		Name:         "Alexandros Papadopoulos",
		Tier:         "Preferred",
		InvoiceTerms: false,
	}

	products := []engine.ProductFact{
		{
			ID:                "PRD-001",
			Name:              "Standard Workstation",
			Category:          "Standard",
			UnitPriceCents:    65000,
			AvailableQuantity: 12,
		},
		{
			ID:                "PRD-002",
			Name:              "Configured Workstation",
			Category:          "CustomBuild",
			UnitPriceCents:    125000,
			AvailableQuantity: 2,
		},
	}

	quote := engine.QuoteFact{
		ID:              "Q-1001",
		CustomerID:      customer.ID,
		DiscountPercent: 20,
		Status:          "Draft",
		Lines: []engine.QuoteLineFact{
			{
				ProductID:      "PRD-002",
				Quantity:       1,
				UnitPriceCents: 125000,
			},
		},
	}

	workingMemory := engine.NewWorkingMemory(customer, quote, products)
	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.DiscountRejectionRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))
	ruleEngine.Register(rules.ApprovalWorkflowGateRule{})
	if *disableCustomBuild {
		if !ruleEngine.SetRuleEnabled("Custom Build Approval Rule", false) {
			panic("Custom Build Approval Rule was not registered")
		}
		fmt.Println("Configuration: CustomBuild approval Rule disabled")
	}
	fmt.Println("Executing registered Rules")
	decision, cycles, err := ruleEngine.DecideUntilStable(workingMemory, 5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inference cycles: %d\n", cycles)

	fmt.Println("=== Working Memory ===")
	fmt.Printf("Customer: %s (%s, tier %s)\n",
		workingMemory.Customer.Name,
		workingMemory.Customer.ID,
		workingMemory.Customer.Tier,
	)
	fmt.Printf("Quote: %s for customer %s, status %s, discount %d%%\n",
		workingMemory.Quote.ID,
		workingMemory.Quote.CustomerID,
		workingMemory.Quote.Status,
		workingMemory.Quote.DiscountPercent,
	)

	for _, line := range workingMemory.Quote.Lines {
		fmt.Printf("- line: product %s, quantity %d, unit price %s\n",
			line.ProductID,
			line.Quantity,
			formatCents(line.UnitPriceCents),
		)
	}

	fmt.Println("Products:")
	for _, product := range workingMemory.Products {
		fmt.Printf("- %s, category %s, price %s, stock %d\n",
			product.ID,
			product.Category,
			formatCents(product.UnitPriceCents),
			product.AvailableQuantity,
		)
	}

	fmt.Printf("Rule findings: %d\n", len(workingMemory.Findings))
	for _, finding := range workingMemory.Findings {
		fmt.Printf("- [%s] %s: %s\n", finding.Severity, finding.RuleName, finding.Message)
	}
	fmt.Printf("Policy decision: %s\n", decision.Outcome)
	fmt.Println("Required reviews:")
	for _, review := range decision.RequiredReviews {
		fmt.Printf("- %s\n", review)
	}
	fmt.Println("Decision reasons:")
	for _, reason := range decision.Reasons {
		fmt.Printf("- %s\n", reason.Message)
	}
	fmt.Println("Derived facts:")
	for _, fact := range workingMemory.DerivedFacts {
		fmt.Printf("- %s = %s\n", fact.Name, fact.Value)
	}
	fmt.Println("Rule trace:")
	for _, trace := range workingMemory.Trace {
		fmt.Printf(
			"- cycle %d, %s: evaluated=%t matched=%t executed=%t",
			trace.Cycle,
			trace.RuleName,
			trace.Evaluated,
			trace.Matched,
			trace.Executed,
		)
		if trace.SkippedReason != "" {
			fmt.Printf(" (%s)", trace.SkippedReason)
		}
		fmt.Println()
	}

	if *simulateQuoteEdit {
		fmt.Println("\n=== Simulated Quote Edit ===")
		fmt.Println("Quote edit: PRD-002 (CustomBuild) -> PRD-001 (Standard), discount 20% -> 0%")
		workingMemory.Quote.Lines = []engine.QuoteLineFact{
			{
				ProductID:      "PRD-001",
				Quantity:       1,
				UnitPriceCents: 65000,
			},
		}
		workingMemory.Quote.DiscountPercent = 0

		decision, cycles, err = ruleEngine.RecomputeDecision(workingMemory, 5)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Inference cycles after recomputation: %d\n", cycles)
		fmt.Printf("Policy decision after recomputation: %s\n", decision.Outcome)
		fmt.Printf("Findings after recomputation: %d\n", len(workingMemory.Findings))
		fmt.Printf("Derived facts after recomputation: %d\n", len(workingMemory.DerivedFacts))
	}
}
