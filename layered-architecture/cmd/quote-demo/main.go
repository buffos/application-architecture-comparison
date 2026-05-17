package main

import (
	"fmt"
	"log"

	"layered-architecture/internal/application"
	"layered-architecture/internal/infrastructure/memory"
	"layered-architecture/internal/presentation/console"
)

func main() {
	customerRepo := memory.NewCustomerRepository()
	productRepo := memory.NewProductRepository()
	quoteRepo := memory.NewQuoteRepository()
	stockRepo := memory.NewStockRecordRepository()
	orderRepo := memory.NewOrderRepository()
	shipmentRepo := memory.NewShipmentRepository()

	customerService := application.NewCustomerService(customerRepo)
	catalogService := application.NewCatalogService(productRepo)
	inventoryService := application.NewInventoryService(productRepo, stockRepo)
	quoteService := application.NewQuoteService(quoteRepo, customerRepo, productRepo)
	orderService := application.NewOrderService(orderRepo, quoteRepo, stockRepo)
	paymentService := application.NewPaymentService(orderRepo)
	fulfillmentService := application.NewFulfillmentService(orderRepo, stockRepo, shipmentRepo)
	handler := console.NewQuoteHandler(customerService, catalogService, inventoryService, quoteService, orderService, paymentService, fulfillmentService)

	output, err := handler.RunDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
