package main

import (
	"fmt"
	"log"

	"modular-monolith/internal/modules/approvals"
	"modular-monolith/internal/modules/customers"
	"modular-monolith/internal/modules/idempotency"
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/products"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/modules/reporting"
	"modular-monolith/internal/modules/returneligibility"
	"modular-monolith/internal/modules/returns"
	"modular-monolith/internal/modules/shipments"
	"modular-monolith/internal/platform/memory"
	paymentadapter "modular-monolith/internal/platform/services/payment"
	timeadapter "modular-monolith/internal/platform/time"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	inventoryRepository := memory.NewInventoryRepository()
	orderRepository := memory.NewOrderRepository()
	productRepository := memory.NewProductRepository()
	quoteRepository := memory.NewQuoteRepository()
	returnRequestRepository := memory.NewReturnRequestRepository()
	shipmentRepository := memory.NewShipmentRepository()
	idempotencyStore := memory.NewIdempotencyStore()

	if err := customerRepository.Save(customers.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(products.Product{
		SKU:              "sku-001",
		Name:             "Desk",
		Category:         "Standard",
		Active:           true,
		UnitPrice:        15000,
		ReturnWindowDays: 30,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(products.Product{
		SKU:              "sku-002",
		Name:             "Custom Desk",
		Category:         "CustomBuild",
		Active:           true,
		UnitPrice:        45000,
		ReturnWindowDays: 14,
	}); err != nil {
		log.Fatal(err)
	}

	if err := inventoryRepository.Save(inventory.StockRecord{
		ProductSKU: "sku-001",
		Available:  10,
	}); err != nil {
		log.Fatal(err)
	}

	if err := inventoryRepository.Save(inventory.StockRecord{
		ProductSKU: "sku-002",
		Available:  3,
	}); err != nil {
		log.Fatal(err)
	}

	customerModule := customers.NewService(customerRepository)
	inventoryModule := inventory.NewService(inventoryRepository)
	paymentModule := payments.NewService(paymentadapter.NewManualReviewGateway())
	productModule := products.NewService(productRepository)
	approvalModule := approvals.NewService()
	clock := timeadapter.NewSystemClock()
	idempotencyModule := idempotency.NewService(idempotencyStore)
	quoteModule := quotes.NewService(quoteRepository, customerModule, productModule, approvalModule)
	returnEligibilityModule := returneligibility.NewService()
	shipmentModule := shipments.NewService(shipmentRepository)
	orderModule := orders.NewService(orderRepository, quoteModule, inventoryModule, paymentModule, shipmentModule, clock)
	returnModule := returns.NewService(returnRequestRepository, orderModule, returnEligibilityModule, inventoryModule, idempotencyModule, paymentModule, clock)
	reportingModule := reporting.NewService(quoteModule, orderModule, returnModule, inventoryModule)

	result, err := quoteModule.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := quoteModule.AddQuoteLine(quotes.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	submitResult, err := quoteModule.SubmitQuote(quotes.SubmitQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted quote: id=%s lines=%d items=%d status=%s\n", submitResult.QuoteID, submitResult.LineCount, submitResult.TotalItems, submitResult.Status)

	details, err := quoteModule.GetQuote(quotes.GetQuoteQuery{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)

	quoteList, err := quoteModule.ListQuotes(quotes.ListQuotesQuery{
		Status: quotes.QuoteStatusApproved,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listed approved quotes: count=%d\n", len(quoteList))

	productDetails, err := productModule.GetProduct(products.GetProductQuery{
		SKU: "sku-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded product: sku=%s category=%s active=%t price=%d\n", productDetails.SKU, productDetails.Category, productDetails.Active, productDetails.UnitPrice)

	productList, err := productModule.ListProducts(products.ListProductsQuery{
		Category:   "Standard",
		ActiveOnly: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listed active standard products: count=%d\n", len(productList))

	customerDetails, err := customerModule.GetCustomer(customers.GetCustomerQuery{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded customer: id=%s active=%t\n", customerDetails.CustomerID, customerDetails.Active)

	customerList, err := customerModule.ListCustomers(customers.ListCustomersQuery{
		ActiveOnly: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listed active customers: count=%d\n", len(customerList))

	pendingResult, err := quoteModule.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = quoteModule.AddQuoteLine(quotes.AddQuoteLineCommand{
		QuoteID:    pendingResult.QuoteID,
		ProductSKU: "sku-002",
		Quantity:   2,
	})
	if err != nil {
		log.Fatal(err)
	}

	pendingSubmit, err := quoteModule.SubmitQuote(quotes.SubmitQuoteCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted custom quote: id=%s lines=%d items=%d status=%s\n", pendingSubmit.QuoteID, pendingSubmit.LineCount, pendingSubmit.TotalItems, pendingSubmit.Status)

	approvalQueue, err := reportingModule.OrdersAwaitingApprovalReport()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range approvalQueue.Rows {
		fmt.Printf("approval queue item: quote=%s customer=%s lines=%d total=%d\n", row.QuoteID, row.CustomerID, row.LineCount, row.TotalAmount)
	}

	approvedPending, err := quoteModule.ApproveQuote(quotes.ApproveQuoteCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("approved pending quote: id=%s lines=%d items=%d status=%s\n", approvedPending.QuoteID, approvedPending.LineCount, approvedPending.TotalItems, approvedPending.Status)

	orderResult, err := orderModule.ConvertQuoteToOrder(orders.ConvertQuoteToOrderCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("converted quote to order: order=%s quote=%s customer=%s status=%s lines=%d\n", orderResult.OrderID, orderResult.QuoteID, orderResult.CustomerID, orderResult.Status, orderResult.LineCount)

	conversionReport, err := reportingModule.QuoteConversionReport()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("quote conversion report: total=%d approved=%d converted=%d rate=%.2f\n", conversionReport.TotalQuotes, conversionReport.ApprovedQuotes, conversionReport.ConvertedQuotes, conversionReport.ConversionRate)

	paidResult, err := orderModule.CapturePayment(orders.CapturePaymentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("captured payment: order=%s customer=%s status=%s lines=%d\n", paidResult.OrderID, paidResult.CustomerID, paidResult.Status, paidResult.LineCount)

	approvedPayment, err := orderModule.ApprovePaymentReview(orders.ApprovePaymentReviewCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("approved payment review: order=%s customer=%s status=%s lines=%d\n", approvedPayment.OrderID, approvedPayment.CustomerID, approvedPayment.Status, approvedPayment.LineCount)

	shipmentResult, err := orderModule.CreateShipment(orders.CreateShipmentCommand{
		OrderID: orderResult.OrderID,
		Lines: []orders.CreateShipmentLine{
			{ProductSKU: "sku-002", Quantity: 1},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created shipment: shipment=%s order=%s customer=%s status=%s lines=%d\n", shipmentResult.ShipmentID, shipmentResult.OrderID, shipmentResult.CustomerID, shipmentResult.Status, shipmentResult.LineCount)

	finalShipmentResult, err := orderModule.CreateShipment(orders.CreateShipmentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created final shipment: shipment=%s order=%s customer=%s status=%s lines=%d\n", finalShipmentResult.ShipmentID, finalShipmentResult.OrderID, finalShipmentResult.CustomerID, finalShipmentResult.Status, finalShipmentResult.LineCount)

	orderDetails, err := orderModule.GetOrder(orders.GetOrderQuery{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded order: order=%s quote=%s status=%s lines=%d\n", orderDetails.OrderID, orderDetails.QuoteID, orderDetails.Status, orderDetails.LineCount)

	orderList, err := orderModule.ListOrders(orders.ListOrdersQuery{
		Status: orders.OrderStatusShipped,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listed shipped orders: count=%d\n", len(orderList))

	shipmentDetails, err := shipmentModule.GetShipment(shipments.GetShipmentQuery{
		ShipmentID: shipmentResult.ShipmentID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded shipment: shipment=%s order=%s customer=%s lines=%d\n", shipmentDetails.ShipmentID, shipmentDetails.OrderID, shipmentDetails.CustomerID, shipmentDetails.LineCount)

	shipmentList, err := shipmentModule.ListShipments(shipments.ListShipmentsQuery{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listed order shipments: count=%d\n", len(shipmentList))

	returnResult, err := returnModule.RequestReturn(returns.RequestReturnCommand{
		OrderID:     orderResult.OrderID,
		Reason:      "damaged item",
		RequestedBy: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("requested return: return=%s order=%s customer=%s status=%s lines=%d\n", returnResult.ReturnRequestID, returnResult.OrderID, returnResult.CustomerID, returnResult.Status, returnResult.LineCount)

	acceptedReturn, err := returnModule.AcceptReturn(returns.ReviewReturnCommand{
		ReturnRequestID: returnResult.ReturnRequestID,
		IdempotencyKey:  "accept-return-001",
		ActorID:         "agent-001",
		ReviewNote:      "accepted after inspection",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("accepted return: return=%s order=%s customer=%s status=%s lines=%d\n", acceptedReturn.ReturnRequestID, acceptedReturn.OrderID, acceptedReturn.CustomerID, acceptedReturn.Status, acceptedReturn.LineCount)

	returnRateReport, err := reportingModule.ReturnRateByCategoryReport()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range returnRateReport.Rows {
		fmt.Printf("return rate by category: category=%s shipped=%d returned=%d rate=%.2f\n", row.Category, row.ShippedQuantity, row.ReturnedQuantity, row.ReturnRate)
	}

	lowStockReport, err := reportingModule.LowStockItemsReport(5)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range lowStockReport.Rows {
		fmt.Printf("low stock item: sku=%s available=%d\n", row.ProductSKU, row.Available)
	}

	returnDetails, err := returnModule.GetReturnRequest(returns.GetReturnRequestQuery{
		ReturnRequestID: returnResult.ReturnRequestID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded return: return=%s order=%s status=%s requestedBy=%s reviewedBy=%s\n", returnDetails.ReturnRequestID, returnDetails.OrderID, returnDetails.Status, returnDetails.RequestedBy, returnDetails.ReviewedBy)

	returnList, err := returnModule.ListReturnRequests(returns.ListReturnRequestsQuery{
		Status: returns.ReturnRequestStatusRefunded,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("listed refunded returns: count=%d\n", len(returnList))
}
