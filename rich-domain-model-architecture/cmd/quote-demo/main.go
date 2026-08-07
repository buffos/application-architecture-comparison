package main

import (
	"fmt"
	"log"
	"time"

	"rich-domain-model-architecture/internal/domain/catalog"
	"rich-domain-model-architecture/internal/domain/customer"
	"rich-domain-model-architecture/internal/domain/fulfillment"
	"rich-domain-model-architecture/internal/domain/inventory"
	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/payments"
	"rich-domain-model-architecture/internal/domain/quoting"
	domainreturns "rich-domain-model-architecture/internal/domain/returns"
	applicationreturns "rich-domain-model-architecture/internal/application/returns"
)

func main() {
	customerAggregate, err := customer.NewCustomer("customer-001", customer.CustomerTierPreferred, customer.PaymentTermsInvoice30)
	if err != nil {
		log.Fatal(err)
	}
	if err := customerAggregate.EnsureCanCreateQuote(); err != nil {
		log.Fatal(err)
	}

	quote, err := quoting.NewQuote("quote-001", quoting.CustomerID(customerAggregate.ID()))
	if err != nil {
		log.Fatal(err)
	}

	catalogPrice, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Desk", catalog.ProductCategoryStandard, catalogPrice)
	if err != nil {
		log.Fatal(err)
	}
	if err := product.EnsureSellable(); err != nil {
		log.Fatal(err)
	}

	quotePrice, err := quoting.NewMoney(product.BasePrice().Cents(), product.BasePrice().Currency())
	if err != nil {
		log.Fatal(err)
	}
	line, err := quoting.NewQuotePricingService().PriceLine(quoting.ProductPricingInput{
		SKU:         quoting.ProductSKU(product.SKU()),
		ProductName: product.Name(),
		Category:    quoting.ProductCategory(product.Category()),
		BasePrice:   quotePrice,
	}, quoting.CustomerPricingTier(customerAggregate.Tier()), 2)
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		log.Fatal(err)
	}

	newCatalogPrice, err := catalog.NewPrice(18000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	if err := product.ChangePrice(newCatalogPrice); err != nil {
		log.Fatal(err)
	}
	approvalDecision := quoting.NewQuoteApprovalService().Evaluate(quote)

	total, err := quote.Total()
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.SubmitForApproval(approvalDecision); err != nil {
		log.Fatal(err)
	}

	snapshot := quote.Lines()[0]
	fmt.Printf("customer domain object: id=%s tier=%s terms=%s active=%t\n", customerAggregate.ID(), customerAggregate.Tier(), customerAggregate.PaymentTerms(), customerAggregate.Active())
	fmt.Printf("product domain object: sku=%s active=%t price=%d %s\n", product.SKU(), product.Active(), product.BasePrice().Cents(), product.BasePrice().Currency())
	fmt.Printf("quote snapshot: sku=%s name=%s unit-price=%d %s\n", snapshot.ProductSKU(), snapshot.ProductName(), snapshot.UnitPrice().Cents(), snapshot.UnitPrice().Currency())
	fmt.Printf("approval decision: required=%t reasons=%d\n", approvalDecision.Required, len(approvalDecision.Reasons))
	fmt.Printf("domain aggregate: id=%s customer=%s status=%s lines=%d total=%d %s\n", quote.ID(), quote.CustomerID(), quote.Status(), len(quote.Lines()), total.Cents(), total.Currency())

	if approvalDecision.Required {
		if err := quote.Approve(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("approved aggregate: id=%s status=%s\n", quote.ID(), quote.Status())

	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		log.Fatal(err)
	}
	orderTotal, err := order.Total()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("order aggregate: id=%s quote=%s status=%s total=%d %s\n", order.ID(), order.QuoteID(), order.Status(), orderTotal.Cents(), orderTotal.Currency())

	stock, err := inventory.NewStockRecord(inventory.ProductSKU(product.SKU()), 10, 2)
	if err != nil {
		log.Fatal(err)
	}
	reservations, err := inventory.NewInventoryReservationService().ReserveAll(
		map[inventory.ProductSKU]*inventory.StockRecord{stock.SKU(): &stock},
		[]inventory.ReservationRequest{{SKU: stock.SKU(), Quantity: order.Lines()[0].Quantity()}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inventory reservation: sku=%s quantity=%d available=%d\n", reservations[0].SKU, reservations[0].Quantity, stock.Available())

	paymentAmount, err := payments.NewMoney(orderTotal.Cents(), orderTotal.Currency())
	if err != nil {
		log.Fatal(err)
	}
	payment, err := payments.NewPayment("payment-001", payments.OrderID(order.ID()), paymentAmount)
	if err != nil {
		log.Fatal(err)
	}
	if err := payment.Capture(); err != nil {
		log.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("payment aggregate: id=%s status=%s order=%s\n", payment.ID(), payment.Status(), payment.OrderID())
	fmt.Printf("paid order aggregate: id=%s status=%s\n", order.ID(), order.Status())

	shipment, err := fulfillment.NewShipmentFromPaidOrder("shipment-001", order)
	if err != nil {
		log.Fatal(err)
	}
	if err := shipment.Dispatch(); err != nil {
		log.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("shipment aggregate: id=%s status=%s lines=%d\n", shipment.ID(), shipment.Status(), len(shipment.Lines()))
	fmt.Printf("shipped order aggregate: id=%s status=%s\n", order.ID(), order.Status())
	if err := order.Cancel(); err != nil {
		fmt.Printf("cancellation blocked: order=%s reason=%s\n", order.ID(), err)
	}

	returnRequest, err := domainreturns.NewReturnRequestFromShippedOrder("return-001", order, "damaged")
	if err != nil {
		log.Fatal(err)
	}
	refundAmount, err := domainreturns.NewMoney(orderTotal.Cents(), orderTotal.Currency())
	if err != nil {
		log.Fatal(err)
	}
	refund, err := domainreturns.NewRefund("refund-001", returnRequest.ID(), refundAmount)
	if err != nil {
		log.Fatal(err)
	}
	if err := refund.Issue(); err != nil {
		log.Fatal(err)
	}
	eligibilityLines := make([]domainreturns.EligibilityLine, 0, len(returnRequest.Lines()))
	for _, line := range returnRequest.Lines() {
		eligibilityLines = append(eligibilityLines, domainreturns.EligibilityLine{Category: line.ProductCategory(), ReturnWindowDays: product.ReturnWindowDays()})
	}
	eligibilityDecision := domainreturns.NewReturnEligibilityService().Evaluate(eligibilityLines)
	windowDecision := domainreturns.NewReturnEligibilityService().EvaluateWindow(
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		eligibilityLines,
	)
	fmt.Printf("return request: id=%s status=%s lines=%d\n", returnRequest.ID(), returnRequest.Status(), len(returnRequest.Lines()))
	fmt.Printf("refund aggregate: id=%s status=%s amount=%d\n", refund.ID(), refund.Status(), refund.Amount().Cents())
	fmt.Printf("return eligibility: eligible=%t reason=%s\n", eligibilityDecision.Eligible, eligibilityDecision.Reason)
	fmt.Printf("return window: eligible=%t reason=%s\n", windowDecision.Eligible, windowDecision.Reason)
	if err := returnRequest.AssignRequester("customer-001"); err != nil {
		log.Fatal(err)
	}
	reviewService := applicationreturns.NewReviewService(applicationreturns.NewInMemoryIdempotencyStore())
	reviewResult, err := reviewService.Review(&returnRequest, domainreturns.ReviewDecisionAccept, "reviewer-001", "processor-001", "return-review-001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("return review: id=%s status=%s reviewed=%s processed=%s\n", reviewResult.ReturnRequestID, reviewResult.Status, reviewResult.ReviewedBy, reviewResult.ProcessedBy)
	restockRequests := make([]inventory.RestockRequest, 0, len(returnRequest.Lines()))
	for _, line := range returnRequest.Lines() {
		restockRequests = append(restockRequests, inventory.RestockRequest{SKU: inventory.ProductSKU(line.ProductSKU()), Quantity: line.Quantity()})
	}
	if err := inventory.NewReturnRestockingService().RestockAll(
		map[inventory.ProductSKU]*inventory.StockRecord{stock.SKU(): &stock},
		restockRequests,
	); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("return restock: sku=%s on-hand=%d available=%d\n", stock.SKU(), stock.OnHand(), stock.Available())
}
