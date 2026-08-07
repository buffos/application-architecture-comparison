package main

import (
	"fmt"

	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/rules"
)

func formatCents(amountCents int64) string {
	return fmt.Sprintf("%d.%02d", amountCents/100, amountCents%100)
}

func main() {
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
	fmt.Println("Executing registered Rules")
	if err := ruleEngine.ExecuteAll(workingMemory); err != nil {
		panic(err)
	}

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
}
