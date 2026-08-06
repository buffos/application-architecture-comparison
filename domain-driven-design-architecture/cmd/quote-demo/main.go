package main

import (
	"fmt"
	"log"

	"domain-driven-design-architecture/internal/domain/catalog"
	"domain-driven-design-architecture/internal/domain/customer"
	"domain-driven-design-architecture/internal/domain/fulfillment"
	"domain-driven-design-architecture/internal/domain/inventory"
	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/payments"
	"domain-driven-design-architecture/internal/domain/quoting"
	"domain-driven-design-architecture/internal/domain/returns"
)

func main() {
	customerAggregate, err := customer.NewCustomer("customer-001", customer.CustomerTierPreferred, customer.PaymentTermsInvoice30)
	if err != nil {
		log.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-001", quoting.CustomerID(customerAggregate.ID()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("customer aggregate: id=%s tier=%s active=%t\n", customerAggregate.ID(), customerAggregate.Tier(), customerAggregate.Active())
	productPrice, err := catalog.NewPrice(15000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	product, err := catalog.NewProduct("sku-001", "Desk", catalog.ProductCategoryStandard, productPrice, 30)
	if err != nil {
		log.Fatal(err)
	}
	if err := product.EnsureSellable(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("product aggregate: sku=%s category=%s active=%t\n", product.SKU(), product.Category(), product.Active())
	stock, err := inventory.NewStockRecord(inventory.ProductSKU(product.SKU()), 10, 3)
	if err != nil {
		log.Fatal(err)
	}
	if err := stock.Reserve(2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("stock aggregate: sku=%s available=%d reserved=%d low=%t\n", stock.SKU(), stock.Available(), stock.Reserved(), stock.IsLowStock())
	quotePrice, err := quoting.NewMoney(product.BasePrice().Cents(), product.BasePrice().Currency())
	if err != nil {
		log.Fatal(err)
	}
	pricingService := quoting.NewQuotePricingService()
	line, err := pricingService.PriceLine(quoting.ProductPricingInput{SKU: quoting.ProductSKU(product.SKU()), Category: quoting.ProductCategory(product.Category()), BasePrice: quotePrice}, quoting.CustomerPricingTier(customerAggregate.Tier()), 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("priced quote line: sku=%s tier=%s unit-price=%d\n", line.ProductSKU(), customerAggregate.Tier(), line.UnitPrice().Cents())
	if err := quote.AddLine(line); err != nil {
		log.Fatal(err)
	}
	approvalDecision := quoting.NewQuoteApprovalService().Evaluate(quote)
	fmt.Printf("quote approval: required=%t reasons=%d\n", approvalDecision.Required, len(approvalDecision.Reasons))
	total, err := quote.Total()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("quote aggregate: id=%s status=%s lines=%d total=%d %s\n", quote.ID(), quote.Status(), len(quote.Lines()), total.Cents(), total.Currency())
	if err := quote.SubmitForApproval(approvalDecision); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("submitted quote aggregate: id=%s status=%s\n", quote.ID(), quote.Status())
	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		log.Fatal(err)
	}
	orderTotal, err := order.Total()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("order aggregate: id=%s quote=%s status=%s total=%d %s\n", order.ID(), order.QuoteID(), order.Status(), orderTotal.Cents(), orderTotal.Currency())
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
	returnRequest, err := returns.NewReturnRequestFromShippedOrder("return-001", order, "damaged")
	if err != nil {
		log.Fatal(err)
	}
	if err := returnRequest.Review(returns.ReviewDecisionAccept); err != nil {
		log.Fatal(err)
	}
	eligibilityLines := make([]returns.EligibilityLine, 0, len(returnRequest.Lines()))
	for _, line := range returnRequest.Lines() {
		eligibilityLines = append(eligibilityLines, returns.EligibilityLine{Category: returns.ProductCategory(line.ProductCategory())})
	}
	eligibilityDecision := returns.NewReturnEligibilityService().Evaluate(eligibilityLines)
	fmt.Printf("return eligibility: eligible=%t reason=%s\n", eligibilityDecision.Eligible, eligibilityDecision.Reason)
	refundAmount, err := returns.NewMoney(orderTotal.Cents(), orderTotal.Currency())
	if err != nil {
		log.Fatal(err)
	}
	refund, err := returns.NewRefund("refund-001", returnRequest.ID(), refundAmount)
	if err != nil {
		log.Fatal(err)
	}
	if err := refund.Issue(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("return request: id=%s status=%s lines=%d\n", returnRequest.ID(), returnRequest.Status(), len(returnRequest.Lines()))
	fmt.Printf("refund aggregate: id=%s status=%s amount=%d\n", refund.ID(), refund.Status(), refund.Amount().Cents())
	stockRecords := map[inventory.ProductSKU]*inventory.StockRecord{stock.SKU(): &stock}
	restockRequests := make([]inventory.RestockRequest, 0, len(returnRequest.Lines()))
	for _, line := range returnRequest.Lines() {
		restockRequests = append(restockRequests, inventory.RestockRequest{SKU: inventory.ProductSKU(line.ProductSKU()), Quantity: line.Quantity()})
	}
	if err := inventory.NewReturnRestockingService().RestockAll(stockRecords, restockRequests); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("return restock: sku=%s on-hand=%d available=%d\n", stock.SKU(), stock.OnHand(), stock.Available())
	reservations, err := inventory.NewInventoryReservationService().ReserveAll(stockRecords, []inventory.ReservationRequest{{SKU: stock.SKU(), Quantity: 2}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inventory reservation: sku=%s quantity=%d available=%d\n", reservations[0].SKU, reservations[0].Quantity, stock.Available())
}
