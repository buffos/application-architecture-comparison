package main

import (
	"flag"
	"fmt"

	"rules-engine-architecture/internal/application"
	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/readmodel"
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
	simulateStockShortage := flag.Bool(
		"simulate-stock-shortage",
		false,
		"set the CustomBuild product stock below the requested quantity",
	)
	simulateShipment := flag.Bool(
		"simulate-shipment",
		false,
		"request shipment for the quote",
	)
	simulatePaymentFailure := flag.Bool(
		"simulate-payment-failure",
		false,
		"set the payment status to failed",
	)
	simulateCancellation := flag.Bool(
		"simulate-cancellation",
		false,
		"request cancellation for the order",
	)
	simulateShippedOrder := flag.Bool(
		"simulate-shipped-order",
		false,
		"set the order status to shipped",
	)
	simulateManagerApproval := flag.Bool(
		"simulate-manager-approval",
		false,
		"load an approved manager-approval Fact",
	)
	simulateReturn := flag.Bool(
		"simulate-return",
		false,
		"request a return for the order",
	)
	simulateClearanceReturn := flag.Bool(
		"simulate-clearance-return",
		false,
		"make the simulated return product clearance",
	)
	returnCommandKey := flag.String(
		"return-command-key",
		"",
		"evaluate the return through an idempotent command key",
	)
	simulateReturnRetry := flag.Bool(
		"simulate-return-retry",
		false,
		"send the same return command twice",
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
			ShortagePolicy:    engine.StockShortageBackorder,
		},
	}
	if *simulateStockShortage {
		products[1].AvailableQuantity = 0
		fmt.Println("Configuration: simulated stock shortage for PRD-002")
	}
	if *simulateClearanceReturn {
		products[1].Category = "Clearance"
		fmt.Println("Configuration: simulated return product is Clearance")
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
	paymentStatus := engine.PaymentAccepted
	if *simulatePaymentFailure {
		paymentStatus = engine.PaymentFailed
	}
	workingMemory.Payment = engine.PaymentFact{Status: paymentStatus}
	workingMemory.Shipment = engine.ShipmentRequestFact{Requested: *simulateShipment}
	orderStatus := engine.OrderConfirmed
	if *simulateShippedOrder {
		orderStatus = engine.OrderShipped
	}
	workingMemory.Order = engine.OrderFact{ID: "ORD-1001", Status: orderStatus}
	workingMemory.Cancellation = engine.CancellationRequestFact{Requested: *simulateCancellation}
	if *simulateReturn {
		workingMemory.ReturnRequest = engine.ReturnRequestFact{
			Requested:         true,
			ProductID:         "PRD-002",
			Quantity:          1,
			ShippedQuantity:   1,
			DaysSinceShipment: 5,
			ReturnWindowDays:  30,
			RequestedBy: engine.ActorFact{
				ID:   "warehouse-clerk-001",
				Role: "warehouse-clerk",
			},
		}
	}
	if *simulateManagerApproval {
		workingMemory.ManagerApproval = engine.ApprovalFact{
			Status:     engine.ApprovalApproved,
			ApprovedBy: "demo-manager",
		}
	}
	ruleEngine := engine.NewEngine()
	ruleEngine.Register(rules.DiscountApprovalRule{})
	ruleEngine.Register(rules.DiscountRejectionRule{})
	ruleEngine.Register(rules.CustomBuildApprovalRule{})
	ruleEngine.Register(rules.NewHighValuePaymentReviewRule(100000))
	ruleEngine.Register(rules.PreferredDiscountEligibilityRule{})
	ruleEngine.Register(rules.InventoryShortageRule{})
	ruleEngine.Register(rules.ShipmentPaymentGuardRule{})
	ruleEngine.Register(rules.CancellationGuardRule{})
	ruleEngine.Register(rules.ReturnPolicyRule{})
	ruleEngine.Register(rules.ApprovalWorkflowGateRule{})
	if *disableCustomBuild {
		if !ruleEngine.SetRuleEnabled("Custom Build Approval Rule", false) {
			panic("Custom Build Approval Rule was not registered")
		}
		fmt.Println("Configuration: CustomBuild approval Rule disabled")
	}
	fmt.Println("Executing registered Rules")
	var decision engine.PolicyDecision
	var cycles int
	var err error
	if *returnCommandKey != "" {
		service := application.NewReturnDecisionService(
			ruleEngine,
			application.NewIdempotencyStore(),
		)
		var replayed bool
		decision, cycles, replayed, err = service.Evaluate(*returnCommandKey, workingMemory, 5)
		if err == nil {
			fmt.Printf("Idempotency replayed: %t\n", replayed)
		}
		if err == nil && *simulateReturnRetry {
			var retryDecision engine.PolicyDecision
			var retryCycles int
			var retryReplayed bool
			retryDecision, retryCycles, retryReplayed, err = service.Evaluate(
				*returnCommandKey,
				workingMemory,
				5,
			)
			if err == nil {
				decision = retryDecision
				fmt.Printf("Retry idempotency replayed: %t\n", retryReplayed)
				fmt.Printf("Retry inference cycles: %d\n", retryCycles)
			}
		}
	} else {
		decision, cycles, err = ruleEngine.DecideUntilStable(workingMemory, 5)
	}
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
	fmt.Printf("Manager approval: %s\n", workingMemory.ManagerApproval.Status)
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
	fmt.Printf("Fulfillment action: %s\n", decision.FulfillmentAction)
	fmt.Printf("Shipment action: %s\n", decision.ShipmentAction)
	fmt.Printf("Cancellation action: %s\n", decision.CancellationAction)
	fmt.Printf("Return action: %s\n", decision.ReturnAction)
	if workingMemory.ReturnRequest.Requested {
		returnView := readmodel.ProjectReturn(workingMemory, decision)
		fmt.Printf(
			"Return query view: order=%s product=%s action=%s requester=%s remaining=%d reason=%s\n",
			returnView.OrderID,
			returnView.ProductID,
			returnView.Action,
			returnView.RequesterID,
			returnView.RemainingQuantity,
			returnView.Reason,
		)
	}
	orderView := readmodel.ProjectOrder(workingMemory, decision)
	fmt.Printf(
		"Order query view: id=%s status=%s payment=%s shipment=%s cancellation=%s return=%s outcome=%s\n",
		orderView.ID,
		orderView.Status,
		orderView.PaymentStatus,
		orderView.ShipmentAction,
		orderView.CancellationAction,
		orderView.ReturnAction,
		orderView.Outcome,
	)
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
		for _, fact := range workingMemory.DerivedFacts {
			fmt.Printf("- fresh fact: %s = %s\n", fact.Name, fact.Value)
		}
	}
}
