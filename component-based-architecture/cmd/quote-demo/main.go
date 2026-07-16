package main

import (
	"fmt"
	"log"

	paymentadapter "component-based-architecture/internal/adapters/payments"
	"component-based-architecture/internal/components/approvals"
	"component-based-architecture/internal/components/clock"
	"component-based-architecture/internal/components/customers"
	"component-based-architecture/internal/components/idempotency"
	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/payments"
	"component-based-architecture/internal/components/products"
	"component-based-architecture/internal/components/quotes"
	"component-based-architecture/internal/components/returneligibility"
	"component-based-architecture/internal/components/returns"
	"component-based-architecture/internal/components/shipments"
)

func main() {
	customerComponent := customers.NewComponent()
	if err := customerComponent.Register(customers.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}
	productComponent := products.NewComponent()
	if err := productComponent.Register(products.Product{
		SKU: "sku-001", Name: "Desk", Category: "Standard", Active: true, UnitPrice: 15000, ReturnWindowDays: 30,
	}); err != nil {
		log.Fatal(err)
	}
	if err := productComponent.Register(products.Product{
		SKU: "sku-002", Name: "Custom Desk", Category: "CustomBuild", Active: true, UnitPrice: 45000, ReturnWindowDays: 30,
	}); err != nil {
		log.Fatal(err)
	}

	approvalComponent := approvals.NewComponent()
	quoteComponent := quotes.NewComponent(customerComponent, productComponent, approvalComponent)
	inventoryComponent := inventory.NewComponent()
	inventoryComponent.RegisterStock(inventory.StockRecord{ProductSKU: "sku-001", Available: 10})
	inventoryComponent.RegisterStock(inventory.StockRecord{ProductSKU: "sku-002", Available: 3})
	paymentComponent := payments.NewComponent(paymentadapter.NewAcceptAllGateway())
	clockComponent := clock.NewComponent()
	shipmentComponent := shipments.NewComponent(clockComponent)
	orderComponent := orders.NewComponent(quoteComponent, inventoryComponent, paymentComponent, shipmentComponent)
	returnEligibilityComponent := returneligibility.NewComponent()
	idempotencyComponent := idempotency.NewComponent()
	returnComponent := returns.NewComponent(orderComponent, paymentComponent, inventoryComponent, returnEligibilityComponent, clockComponent, idempotencyComponent)
	result, err := quoteComponent.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := quoteComponent.AddQuoteLine(quotes.AddQuoteLineCommand{
		QuoteID: result.QuoteID, ProductSKU: "sku-001", Quantity: 2,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("added quote line: id=%s lines=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.Status)

	submission, err := quoteComponent.SubmitQuote(quotes.SubmitQuoteCommand{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("submitted quote: id=%s lines=%d status=%s\n", submission.QuoteID, submission.LineCount, submission.Status)

	var quoteLookup quotes.QuoteLookup = quoteComponent
	details, err := quoteLookup.GetQuote(quotes.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)

	pending, err := quoteComponent.CreateDraftQuote(quotes.CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := quoteComponent.AddQuoteLine(quotes.AddQuoteLineCommand{QuoteID: pending.QuoteID, ProductSKU: "sku-002", Quantity: 1}); err != nil {
		log.Fatal(err)
	}
	pendingSubmission, err := quoteComponent.SubmitQuote(quotes.SubmitQuoteCommand{QuoteID: pending.QuoteID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("submitted custom quote: id=%s status=%s\n", pendingSubmission.QuoteID, pendingSubmission.Status)

	approval, err := quoteComponent.ApproveQuote(quotes.ApproveQuoteCommand{QuoteID: pending.QuoteID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("approved custom quote: id=%s status=%s\n", approval.QuoteID, approval.Status)

	order, err := orderComponent.ConvertQuoteToOrder(orders.ConvertQuoteToOrderCommand{QuoteID: pending.QuoteID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("converted quote to order: order=%s quote=%s status=%s lines=%d\n", order.OrderID, order.QuoteID, order.Status, order.LineCount)

	paid, err := orderComponent.CapturePayment(orders.CapturePaymentCommand{OrderID: order.OrderID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("captured payment: order=%s status=%s\n", paid.OrderID, paid.Status)

	shipment, err := orderComponent.CreateShipment(orders.CreateShipmentCommand{OrderID: order.OrderID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created shipment: shipment=%s order=%s status=%s\n", shipment.ShipmentID, shipment.OrderID, shipment.Status)
	returnRequest, err := returnComponent.RequestReturn(returns.RequestReturnCommand{OrderID: order.OrderID, Reason: "damaged", RequestedBy: "agent-001"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("requested return: return=%s order=%s status=%s\n", returnRequest.ReturnRequestID, returnRequest.OrderID, returnRequest.Status)
	var returnReader returns.Reader = returnComponent
	returnDetails, err := returnReader.GetReturnRequest(returns.GetReturnRequestQuery{ReturnRequestID: returnRequest.ReturnRequestID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loaded return: return=%s requester=%s status=%s\n", returnDetails.ReturnRequestID, returnDetails.RequestedBy, returnDetails.Status)
	if _, err := returnComponent.AcceptReturn(returns.ReviewReturnCommand{ReturnRequestID: returnRequest.ReturnRequestID, ReviewedBy: "reviewer-001", ProcessedBy: "processor-001", ReviewNote: "eligible", IdempotencyKey: "accept-return-001"}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("accepted return: return=%s status=%s\n", returnRequest.ReturnRequestID, returns.ReturnRequestStatusRefunded)
	refundedReturns := returnReader.ListReturnRequests(returns.ListReturnRequestsQuery{Status: returns.ReturnRequestStatusRefunded})
	fmt.Printf("listed refunded returns: count=%d\n", len(refundedReturns))

	cancellable, err := quoteComponent.CreateDraftQuote(quotes.CreateDraftQuoteCommand{CustomerID: "customer-001"})
	if err != nil {
		log.Fatal(err)
	}
	if _, err := quoteComponent.AddQuoteLine(quotes.AddQuoteLineCommand{QuoteID: cancellable.QuoteID, ProductSKU: "sku-001", Quantity: 1}); err != nil {
		log.Fatal(err)
	}
	if _, err := quoteComponent.SubmitQuote(quotes.SubmitQuoteCommand{QuoteID: cancellable.QuoteID}); err != nil {
		log.Fatal(err)
	}
	cancellableOrder, err := orderComponent.ConvertQuoteToOrder(orders.ConvertQuoteToOrderCommand{QuoteID: cancellable.QuoteID})
	if err != nil {
		log.Fatal(err)
	}
	cancelled, err := orderComponent.CancelOrder(orders.CancelOrderCommand{OrderID: cancellableOrder.OrderID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cancelled order: order=%s status=%s\n", cancelled.OrderID, cancelled.Status)
}
