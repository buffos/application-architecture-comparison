package main

import (
	"flag"
	"fmt"

	"rules-engine-architecture/internal/application"
	"rules-engine-architecture/internal/engine"
	"rules-engine-architecture/internal/plugins"
	"rules-engine-architecture/internal/readmodel"
	"rules-engine-architecture/internal/reporting"
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
	simulatePartialShipment := flag.Bool(
		"simulate-partial-shipment",
		false,
		"request shipment after one unit has already shipped",
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
	simulatePaymentReviewApproved := flag.Bool(
		"simulate-payment-review-approved",
		false,
		"load an approved high-value payment-review Fact",
	)
	simulatePaymentReviewRejected := flag.Bool(
		"simulate-payment-review-rejected",
		false,
		"load a rejected high-value payment-review Fact",
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
	simulatePartialReturn := flag.Bool(
		"simulate-partial-return",
		false,
		"request a return with one accepted and one rejected line",
	)
	enableSeasonalPlugin := flag.Bool(
		"enable-seasonal-plugin",
		false,
		"register the optional seasonal surcharge Rule plugin",
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
	lowStockThreshold := flag.Int(
		"low-stock-threshold",
		3,
		"report products with stock below this quantity",
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
	if *simulatePartialShipment {
		workingMemory.Shipment = engine.ShipmentRequestFact{
			Requested: true,
			Lines: []engine.ShipmentLineFact{{
				ProductID:              "PRD-002",
				OrderedQuantity:        2,
				AlreadyShippedQuantity: 1,
			}},
		}
		fmt.Println("Configuration: simulated partial shipment for PRD-002")
	}
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
	if *simulatePartialReturn {
		workingMemory.ReturnRequest = engine.ReturnRequestFact{
			Requested:         true,
			DaysSinceShipment: 5,
			ReturnWindowDays:  30,
			RequestedBy: engine.ActorFact{
				ID:   "warehouse-clerk-001",
				Role: "warehouse-clerk",
			},
			Lines: []engine.ReturnLineFact{
				{ProductID: "PRD-002", Quantity: 1, ShippedQuantity: 1},
				{ProductID: "PRD-001", Quantity: 2, ShippedQuantity: 1},
			},
		}
		fmt.Println("Configuration: simulated partial return")
	}
	if *simulateManagerApproval {
		workingMemory.ManagerApproval = engine.ApprovalFact{
			Status:     engine.ApprovalApproved,
			ApprovedBy: "demo-manager",
		}
	}
	if *simulatePaymentReviewApproved && *simulatePaymentReviewRejected {
		panic("payment review cannot be both approved and rejected")
	}
	if *simulatePaymentReviewApproved {
		workingMemory.PaymentReview = engine.PaymentReviewFact{
			Status:     engine.PaymentReviewApproved,
			ReviewedBy: "payment-manager",
		}
	}
	if *simulatePaymentReviewRejected {
		workingMemory.PaymentReview = engine.PaymentReviewFact{
			Status:     engine.PaymentReviewRejected,
			ReviewedBy: "payment-manager",
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
	ruleEngine.Register(rules.PartialShipmentRule{})
	ruleEngine.Register(rules.CancellationGuardRule{})
	ruleEngine.Register(rules.ReturnPolicyRule{})
	ruleEngine.Register(rules.ApprovalWorkflowGateRule{})
	if *enableSeasonalPlugin {
		ruleEngine.Register(plugins.NewSeasonalSurchargeRule("CustomBuild", 5))
		fmt.Println("Configuration: seasonal surcharge plugin enabled")
	}
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
	fmt.Printf("Payment review: %s\n", workingMemory.PaymentReview.Status)
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
	fmt.Printf("Pricing adjustment: %s\n", formatCents(decision.PricingAdjustmentCents))
	if workingMemory.ReturnRequest.Requested {
		returnView := readmodel.ProjectReturn(workingMemory, decision)
		fmt.Printf(
			"Return query view: order=%s product=%s lines=%d action=%s partial=%t requester=%s remaining=%d reason=%s\n",
			returnView.OrderID,
			returnView.ProductID,
			returnView.LineCount,
			returnView.Action,
			returnView.Partial,
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
	quoteList := readmodel.ProjectQuoteList([]readmodel.EvaluatedQuote{{
		Memory:   workingMemory,
		Decision: decision,
	}})
	for _, quoteSummary := range quoteList {
		fmt.Printf(
			"Quote list view: id=%s customer=%s status=%s subtotal=%s discount=%d%% outcome=%s\n",
			quoteSummary.ID,
			quoteSummary.CustomerID,
			quoteSummary.Status,
			formatCents(quoteSummary.SubtotalCents),
			quoteSummary.DiscountPercent,
			quoteSummary.Outcome,
		)
	}
	for _, productView := range readmodel.ProjectProducts(workingMemory.Products) {
		fmt.Printf(
			"Product query view: id=%s category=%s price=%s stock=%d shortage-policy=%s\n",
			productView.ID,
			productView.Category,
			formatCents(productView.UnitPriceCents),
			productView.AvailableQuantity,
			productView.ShortagePolicy,
		)
	}
	for _, row := range reporting.BuildLowStockReport(workingMemory.Products, *lowStockThreshold) {
		fmt.Printf(
			"Low-stock report: product=%s category=%s stock=%d threshold=%d\n",
			row.ProductID,
			row.ProductCategory,
			row.AvailableQuantity,
			row.Threshold,
		)
	}
	for _, customerView := range readmodel.ProjectCustomers([]engine.CustomerFact{workingMemory.Customer}) {
		fmt.Printf(
			"Customer query view: id=%s name=%s tier=%s invoice-terms=%t\n",
			customerView.ID,
			customerView.Name,
			customerView.Tier,
			customerView.InvoiceTerms,
		)
	}
	conversionReport := reporting.BuildQuoteConversionReport([]engine.QuoteFact{workingMemory.Quote})
	fmt.Printf(
		"Quote conversion report: converted=%d/%d rate=%.2f%%\n",
		conversionReport.ConvertedQuotes,
		conversionReport.TotalQuotes,
		conversionReport.ConversionRatePercent,
	)
	if workingMemory.ReturnRequest.Requested {
		category := "Unknown"
		productID := workingMemory.ReturnRequest.ProductID
		if productID == "" && len(workingMemory.ReturnRequest.Lines) > 0 {
			productID = workingMemory.ReturnRequest.Lines[0].ProductID
		}
		for _, product := range workingMemory.Products {
			if product.ID == productID {
				category = product.Category
				break
			}
		}
		returnRateRows := reporting.BuildReturnRateByCategory([]reporting.ReturnRecord{{
			ProductCategory: category,
			Accepted:        decision.ReturnAction == engine.ReturnAllowed,
		}})
		for _, row := range returnRateRows {
			fmt.Printf(
				"Return rate report: category=%s accepted=%d/%d rate=%.2f%%\n",
				row.ProductCategory,
				row.AcceptedReturns,
				row.AttemptedReturns,
				row.ReturnRatePercent,
			)
		}
	}
	approvalRows := reporting.BuildOrdersAwaitingApprovalReport([]readmodel.EvaluatedQuote{{
		Memory:   workingMemory,
		Decision: decision,
	}})
	for _, row := range approvalRows {
		fmt.Printf(
			"Approval queue: quote=%s customer=%s reviews=%v reasons=%d\n",
			row.QuoteID,
			row.CustomerID,
			row.RequiredReviews,
			len(row.Reasons),
		)
	}
	if workingMemory.Shipment.Requested {
		shipmentView := readmodel.ProjectShipment(workingMemory, decision)
		fmt.Printf(
			"Shipment query view: order=%s requested=%t action=%s partial=%t payment=%s invoice-terms=%t reason=%s\n",
			shipmentView.OrderID,
			shipmentView.Requested,
			shipmentView.Action,
			shipmentView.Partial,
			shipmentView.PaymentStatus,
			shipmentView.InvoiceTerms,
			shipmentView.Reason,
		)
	}
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
