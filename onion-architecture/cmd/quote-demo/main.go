package main

import (
	"fmt"
	"log"

	"onion-architecture/internal/application"
	"onion-architecture/internal/domain"
	"onion-architecture/internal/infrastructure/memory"
	"onion-architecture/internal/infrastructure/policies/approval"
	returneligibility "onion-architecture/internal/infrastructure/policies/returneligibility"
	"onion-architecture/internal/infrastructure/services/payment"
	timeinfra "onion-architecture/internal/infrastructure/services/time"
)

func main() {
	customerRepository := memory.NewCustomerRepository()
	quoteRepository := memory.NewQuoteRepository()
	productRepository := memory.NewProductRepository()
	orderRepository := memory.NewOrderRepository()
	shipmentRepository := memory.NewShipmentRepository()
	inventoryReservation := memory.NewInventoryReservation()
	clock := timeinfra.NewSystemClock()

	if err := customerRepository.Save(domain.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(domain.Product{
		SKU:              "sku-001",
		Name:             "Desk",
		Category:         "Standard",
		Active:           true,
		UnitPrice:        15000,
		ReturnWindowDays: 30,
	}); err != nil {
		log.Fatal(err)
	}

	if err := productRepository.Save(domain.Product{
		SKU:              "sku-002",
		Name:             "Custom Desk",
		Category:         "CustomBuild",
		Active:           true,
		UnitPrice:        45000,
		ReturnWindowDays: 30,
	}); err != nil {
		log.Fatal(err)
	}

	inventoryReservation.Seed("sku-002", 5)
	inventoryReservation.Seed("sku-001", 20)

	submissionPolicy := approval.NewCategoryPolicy()
	_ = returneligibility.NewWindowPolicy()
	paymentGateway := payment.NewAcceptAllGateway()
	service := application.NewCreateDraftQuoteService(quoteRepository, customerRepository)
	getQuote := application.NewGetQuoteService(quoteRepository)
	addQuoteLine := application.NewAddQuoteLineService(quoteRepository, productRepository)
	submitQuote := application.NewSubmitQuoteService(quoteRepository, submissionPolicy)
	approveQuote := application.NewApproveQuoteService(quoteRepository)
	convertQuote := application.NewConvertQuoteToOrderService(quoteRepository, orderRepository, inventoryReservation)
	capturePayment := application.NewCapturePaymentService(orderRepository, paymentGateway)
	createShipment := application.NewCreateShipmentService(orderRepository, shipmentRepository, clock)
	lowStockItemsReport := application.NewLowStockItemsReportService(inventoryReservation)

	result, err := service.Execute(application.CreateDraftQuoteCommand{
		CustomerID: "customer-001",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created draft quote: id=%s customer=%s status=%s\n", result.QuoteID, result.CustomerID, result.Status)

	lineResult, err := addQuoteLine.Execute(application.AddQuoteLineCommand{
		QuoteID:    result.QuoteID,
		ProductSKU: "sku-002",
		Quantity:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("added quote line: id=%s lines=%d items=%d status=%s\n", lineResult.QuoteID, lineResult.LineCount, lineResult.TotalItems, lineResult.Status)

	submitResult, err := submitQuote.Execute(application.SubmitQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("submitted quote: id=%s lines=%d items=%d status=%s\n", submitResult.QuoteID, submitResult.LineCount, submitResult.TotalItems, submitResult.Status)

	approvalResult, err := approveQuote.Execute(application.ApproveQuoteCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("approved quote: id=%s lines=%d items=%d status=%s\n", approvalResult.QuoteID, approvalResult.LineCount, approvalResult.TotalItems, approvalResult.Status)

	orderResult, err := convertQuote.Execute(application.ConvertQuoteToOrderCommand{
		QuoteID: result.QuoteID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("converted quote: order=%s quote=%s customer=%s status=%s lines=%d\n", orderResult.OrderID, orderResult.QuoteID, orderResult.CustomerID, orderResult.Status, orderResult.LineCount)

	paymentResult, err := capturePayment.Execute(application.CapturePaymentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("captured payment: order=%s quote=%s customer=%s status=%s lines=%d\n", paymentResult.OrderID, paymentResult.QuoteID, paymentResult.CustomerID, paymentResult.Status, paymentResult.LineCount)

	shipmentResult, err := createShipment.Execute(application.CreateShipmentCommand{
		OrderID: orderResult.OrderID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created shipment: shipment=%s order=%s orderStatus=%s lines=%d\n", shipmentResult.ShipmentID, shipmentResult.OrderID, shipmentResult.OrderStatus, shipmentResult.LineCount)

	details, err := getQuote.Execute(application.GetQuoteQuery{QuoteID: result.QuoteID})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loaded quote: id=%s customer=%s status=%s lines=%d\n", details.QuoteID, details.CustomerID, details.Status, details.LineCount)

	lowStockItems, err := lowStockItemsReport.Execute(application.LowStockItemsReportQuery{
		Threshold: 5,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("low stock items: %v\n", lowStockItems)
}
