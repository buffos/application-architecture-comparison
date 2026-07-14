package main

import (
	"fmt"
	"log"

	"microkernel-architecture/internal/kernel"
	"microkernel-architecture/internal/platform/memory"
	"microkernel-architecture/internal/plugins/approvals"
	"microkernel-architecture/internal/plugins/clock"
	"microkernel-architecture/internal/plugins/customers"
	"microkernel-architecture/internal/plugins/idempotency"
	"microkernel-architecture/internal/plugins/inventory"
	"microkernel-architecture/internal/plugins/orders"
	"microkernel-architecture/internal/plugins/payments"
	"microkernel-architecture/internal/plugins/products"
	"microkernel-architecture/internal/plugins/quotes"
	"microkernel-architecture/internal/plugins/returneligibility"
	"microkernel-architecture/internal/plugins/returns"
	"microkernel-architecture/internal/plugins/shipments"
)

func main() {
	host := kernel.NewHost()

	customerRepository := memory.NewCustomerRepository()
	inventoryRepository := memory.NewInventoryRepository()
	idempotencyStore := memory.NewIdempotencyStore()
	orderRepository := memory.NewOrderRepository()
	productRepository := memory.NewProductRepository()
	quoteRepository := memory.NewQuoteRepository()
	returnRequestRepository := memory.NewReturnRequestRepository()
	shipmentRepository := memory.NewShipmentRepository()

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

	if err := host.RegisterPlugin(customers.NewPlugin(customerRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(products.NewPlugin(productRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(approvals.NewPlugin()); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(clock.NewPlugin()); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(quotes.NewPlugin(quoteRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(inventory.NewPlugin(inventoryRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(payments.NewPlugin()); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(shipments.NewPlugin(shipmentRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(orders.NewPlugin(orderRepository)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(returneligibility.NewPlugin()); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(idempotency.NewPlugin(idempotencyStore)); err != nil {
		log.Fatal(err)
	}

	if err := host.RegisterPlugin(returns.NewPlugin(returnRequestRepository)); err != nil {
		log.Fatal(err)
	}

	quoteService, err := host.QuoteService()
	if err != nil {
		log.Fatal(err)
	}

	quoteReader, err := host.QuoteReader()
	if err != nil {
		log.Fatal(err)
	}

	orderService, err := host.OrderService()
	if err != nil {
		log.Fatal(err)
	}

	returnService, err := host.ReturnService()
	if err != nil {
		log.Fatal(err)
	}

	result, err := quoteService.CreateDraftQuote(kernel.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := quoteService.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-001",
		Quantity:   2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	submitResult, err := quoteService.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted quote: id=%s lines=%d items=%d status=%s\n", submitResult.QuoteID, submitResult.LineCount, submitResult.TotalItems, submitResult.Status)

	details, err := quoteReader.GetQuote(kernel.GetQuoteQuery{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d items=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount, details.TotalItems)

	pendingResult, err := quoteService.CreateDraftQuote(kernel.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = quoteService.AddQuoteLine(kernel.AddQuoteLineCommand{
		QuoteID:    pendingResult.QuoteID,
		ProductSKU: "sku-002",
		Quantity:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	pendingSubmit, err := quoteService.SubmitQuote(kernel.SubmitQuoteCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted custom quote: id=%s status=%s\n", pendingSubmit.QuoteID, pendingSubmit.Status)

	approvedPending, err := quoteService.ApproveQuote(kernel.ApproveQuoteCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("approved pending quote: id=%s status=%s\n", approvedPending.QuoteID, approvedPending.Status)

	orderResult, err := orderService.ConvertQuoteToOrder(kernel.ConvertQuoteToOrderCommand{
		QuoteID: pendingResult.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("converted quote to order: order=%s quote=%s customer=%s status=%s lines=%d\n", orderResult.OrderID, orderResult.QuoteID, orderResult.CustomerID, orderResult.Status, orderResult.LineCount)

	paidResult, err := orderService.CapturePayment(kernel.CapturePaymentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("captured payment: order=%s quote=%s customer=%s status=%s lines=%d\n", paidResult.OrderID, paidResult.QuoteID, paidResult.CustomerID, paidResult.Status, paidResult.LineCount)

	shipmentResult, err := orderService.CreateShipment(kernel.CreateShipmentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created shipment: shipment=%s order=%s customer=%s status=%s lines=%d\n", shipmentResult.ShipmentID, shipmentResult.OrderID, shipmentResult.CustomerID, shipmentResult.Status, shipmentResult.LineCount)

	returnResult, err := returnService.RequestReturn(kernel.RequestReturnCommand{
		OrderID:     orderResult.OrderID,
		Reason:      "damaged item",
		RequestedBy: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("requested return: return=%s order=%s customer=%s status=%s lines=%d\n", returnResult.ReturnRequestID, returnResult.OrderID, returnResult.CustomerID, returnResult.Status, returnResult.LineCount)

	acceptedReturn, err := returnService.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: returnResult.ReturnRequestID,
		IdempotencyKey:  "accept-return-001",
		ReviewedBy:      "agent-001",
		ProcessedBy:     "ops-001",
		ReviewNote:      "approved after inspection",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("accepted return: return=%s order=%s customer=%s status=%s lines=%d\n", acceptedReturn.ReturnRequestID, acceptedReturn.OrderID, acceptedReturn.CustomerID, acceptedReturn.Status, acceptedReturn.LineCount)

	restockedItem, err := inventoryRepository.FindBySKU("sku-002")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("restocked inventory: sku=%s available=%d\n", restockedItem.ProductSKU, restockedItem.Available)
}
