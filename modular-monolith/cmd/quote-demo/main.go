package main

import (
	"fmt"
	"log"

	"modular-monolith/internal/modules/approvals"
	"modular-monolith/internal/modules/customers"
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/products"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/modules/returns"
	"modular-monolith/internal/modules/shipments"
	"modular-monolith/internal/platform/memory"
	paymentadapter "modular-monolith/internal/platform/services/payment"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	inventoryRepository := memory.NewInventoryRepository()
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
		SKU:       "sku-001",
		Name:      "Desk",
		Category:  "Standard",
		Active:    true,
		UnitPrice: 15000,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(products.Product{
		SKU:       "sku-002",
		Name:      "Custom Desk",
		Category:  "CustomBuild",
		Active:    true,
		UnitPrice: 45000,
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
	paymentModule := payments.NewService(paymentadapter.NewAcceptAllGateway())
	productModule := products.NewService(productRepository)
	approvalModule := approvals.NewService()
	quoteModule := quotes.NewService(quoteRepository, customerModule, productModule, approvalModule)
	shipmentModule := shipments.NewService(shipmentRepository)
	orderModule := orders.NewService(orderRepository, quoteModule, inventoryModule, paymentModule, shipmentModule)
	returnModule := returns.NewService(returnRequestRepository, orderModule, inventoryModule, paymentModule)

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

	pendingResult, err := quoteModule.CreateDraftQuote(quotes.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = quoteModule.AddQuoteLine(quotes.AddQuoteLineCommand{
		QuoteID:    pendingResult.QuoteID,
		ProductSKU: "sku-002",
		Quantity:   1,
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

	paidResult, err := orderModule.CapturePayment(orders.CapturePaymentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("captured payment: order=%s customer=%s status=%s lines=%d\n", paidResult.OrderID, paidResult.CustomerID, paidResult.Status, paidResult.LineCount)

	shipmentResult, err := orderModule.CreateShipment(orders.CreateShipmentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created shipment: shipment=%s order=%s customer=%s status=%s lines=%d\n", shipmentResult.ShipmentID, shipmentResult.OrderID, shipmentResult.CustomerID, shipmentResult.Status, shipmentResult.LineCount)

	returnResult, err := returnModule.RequestReturn(returns.RequestReturnCommand{
		OrderID: orderResult.OrderID,
		Reason:  "damaged item",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("requested return: return=%s order=%s customer=%s status=%s lines=%d\n", returnResult.ReturnRequestID, returnResult.OrderID, returnResult.CustomerID, returnResult.Status, returnResult.LineCount)

	acceptedReturn, err := returnModule.AcceptReturn(returns.ReviewReturnCommand{
		ReturnRequestID: returnResult.ReturnRequestID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("accepted return: return=%s order=%s customer=%s status=%s lines=%d\n", acceptedReturn.ReturnRequestID, acceptedReturn.OrderID, acceptedReturn.CustomerID, acceptedReturn.Status, acceptedReturn.LineCount)
}
