package main

import (
	"fmt"
	"log"

	"clean-architecture/internal/entities"
	"clean-architecture/internal/infrastructure/memory"
	approvalpolicy "clean-architecture/internal/infrastructure/policies/approval"
	returneligibility "clean-architecture/internal/infrastructure/policies/returneligibility"
	paymentservice "clean-architecture/internal/infrastructure/services/payment"
	timeadapter "clean-architecture/internal/infrastructure/services/time"
	"clean-architecture/internal/interfaceadapters/controllers"
	"clean-architecture/internal/interfaceadapters/presenters"
	"clean-architecture/internal/usecases"
)

func main() {
	customerGateway := memory.NewCustomerGateway()
	quoteGateway := memory.NewQuoteGateway()
	orderGateway := memory.NewOrderGateway()
	shipmentGateway := memory.NewShipmentGateway()
	returnRequestGateway := memory.NewReturnRequestGateway()
	productGateway := memory.NewProductGateway()
	inventoryReservation := memory.NewInventoryReservation(map[string]int{
		"CHAIR-001": 5,
	})
	approvalPolicy := approvalpolicy.NewCategoryPolicy()
	clock := timeadapter.NewSystemClock()
	_ = returneligibility.NewWindowPolicy()
	paymentGateway := paymentservice.NewAcceptAllGateway()
	createPresenter := presenters.NewCreateDraftQuotePresenter()
	getCustomerPresenter := presenters.NewGetCustomerPresenter()
	listCustomersPresenter := presenters.NewListCustomersPresenter()
	quoteConversionReportPresenter := presenters.NewQuoteConversionReportPresenter()
	addLinePresenter := presenters.NewAddQuoteLinePresenter()
	submitPresenter := presenters.NewSubmitQuotePresenter()
	convertPresenter := presenters.NewConvertQuoteToOrderPresenter()
	capturePresenter := presenters.NewCapturePaymentPresenter()
	shipmentPresenter := presenters.NewCreateShipmentPresenter()
	getPresenter := presenters.NewGetQuotePresenter()
	listQuotesPresenter := presenters.NewListQuotesPresenter()
	getProductPresenter := presenters.NewGetProductPresenter()
	listProductsPresenter := presenters.NewListProductsPresenter()
	getOrderPresenter := presenters.NewGetOrderPresenter()
	listOrdersPresenter := presenters.NewListOrdersPresenter()
	getShipmentPresenter := presenters.NewGetShipmentPresenter()
	listShipmentsPresenter := presenters.NewListShipmentsPresenter()
	requestReturnPresenter := presenters.NewRequestReturnPresenter()
	getReturnPresenter := presenters.NewGetReturnRequestPresenter()
	listReturnPresenter := presenters.NewListReturnRequestsPresenter()

	if err := customerGateway.Save(entities.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}
	if err := productGateway.Save(entities.Product{
		SKU:              "CHAIR-001",
		Name:             "Office Chair",
		Category:         "Standard",
		BasePrice:        10000,
		Available:        true,
		ReturnWindowDays: 30,
	}); err != nil {
		log.Fatal(err)
	}

	createInteractor := usecases.NewCreateDraftQuoteInteractor(quoteGateway, customerGateway, createPresenter)
	createController := controllers.NewCreateDraftQuoteController(createInteractor)

	getCustomerInteractor := usecases.NewGetCustomerInteractor(customerGateway, getCustomerPresenter)
	getCustomerController := controllers.NewGetCustomerController(getCustomerInteractor)

	listCustomersInteractor := usecases.NewListCustomersInteractor(customerGateway, listCustomersPresenter)
	listCustomersController := controllers.NewListCustomersController(listCustomersInteractor)

	quoteConversionReportInteractor := usecases.NewQuoteConversionReportInteractor(quoteGateway, orderGateway, quoteConversionReportPresenter)
	quoteConversionReportController := controllers.NewQuoteConversionReportController(quoteConversionReportInteractor)

	if err := createController.Handle("customer-001"); err != nil {
		log.Fatal(err)
	}

	if err := getCustomerController.Handle("customer-001"); err != nil {
		log.Fatal(err)
	}

	if err := listCustomersController.Handle(true); err != nil {
		log.Fatal(err)
	}

	addLineInteractor := usecases.NewAddQuoteLineInteractor(quoteGateway, productGateway, addLinePresenter)
	addLineController := controllers.NewAddQuoteLineController(addLineInteractor)

	if err := addLineController.Handle(createPresenter.ViewModel().QuoteID, "CHAIR-001", 2); err != nil {
		log.Fatal(err)
	}

	submitInteractor := usecases.NewSubmitQuoteInteractor(quoteGateway, approvalPolicy, submitPresenter)
	submitController := controllers.NewSubmitQuoteController(submitInteractor)

	if err := submitController.Handle(createPresenter.ViewModel().QuoteID); err != nil {
		log.Fatal(err)
	}

	convertInteractor := usecases.NewConvertQuoteToOrderInteractor(quoteGateway, orderGateway, inventoryReservation, clock, convertPresenter)
	convertController := controllers.NewConvertQuoteToOrderController(convertInteractor)

	if err := convertController.Handle(createPresenter.ViewModel().QuoteID); err != nil {
		log.Fatal(err)
	}

	captureInteractor := usecases.NewCapturePaymentInteractor(orderGateway, paymentGateway, capturePresenter)
	captureController := controllers.NewCapturePaymentController(captureInteractor)

	if err := captureController.Handle(convertPresenter.ViewModel().OrderID); err != nil {
		log.Fatal(err)
	}

	createShipmentInteractor := usecases.NewCreateShipmentInteractor(orderGateway, shipmentGateway, clock, shipmentPresenter)
	createShipmentController := controllers.NewCreateShipmentController(createShipmentInteractor)

	if err := createShipmentController.Handle(convertPresenter.ViewModel().OrderID); err != nil {
		log.Fatal(err)
	}

	getInteractor := usecases.NewGetQuoteInteractor(quoteGateway, getPresenter)
	getController := controllers.NewGetQuoteController(getInteractor)

	listQuotesInteractor := usecases.NewListQuotesInteractor(quoteGateway, listQuotesPresenter)
	listQuotesController := controllers.NewListQuotesController(listQuotesInteractor)

	getProductInteractor := usecases.NewGetProductInteractor(productGateway, getProductPresenter)
	getProductController := controllers.NewGetProductController(getProductInteractor)

	listProductsInteractor := usecases.NewListProductsInteractor(productGateway, listProductsPresenter)
	listProductsController := controllers.NewListProductsController(listProductsInteractor)

	getOrderInteractor := usecases.NewGetOrderInteractor(orderGateway, getOrderPresenter)
	getOrderController := controllers.NewGetOrderController(getOrderInteractor)

	listOrdersInteractor := usecases.NewListOrdersInteractor(orderGateway, listOrdersPresenter)
	listOrdersController := controllers.NewListOrdersController(listOrdersInteractor)

	getShipmentInteractor := usecases.NewGetShipmentInteractor(shipmentGateway, getShipmentPresenter)
	getShipmentController := controllers.NewGetShipmentController(getShipmentInteractor)

	listShipmentsInteractor := usecases.NewListShipmentsInteractor(shipmentGateway, listShipmentsPresenter)
	listShipmentsController := controllers.NewListShipmentsController(listShipmentsInteractor)

	requestReturnInteractor := usecases.NewRequestReturnInteractor(orderGateway, returnRequestGateway, clock, requestReturnPresenter)
	requestReturnController := controllers.NewRequestReturnController(requestReturnInteractor)

	getReturnInteractor := usecases.NewGetReturnRequestInteractor(returnRequestGateway, getReturnPresenter)
	getReturnController := controllers.NewGetReturnRequestController(getReturnInteractor)

	listReturnInteractor := usecases.NewListReturnRequestsInteractor(returnRequestGateway, listReturnPresenter)
	listReturnController := controllers.NewListReturnRequestsController(listReturnInteractor)

	if err := getController.Handle(createPresenter.ViewModel().QuoteID); err != nil {
		log.Fatal(err)
	}

	if err := listQuotesController.Handle(entities.QuoteStatusApproved); err != nil {
		log.Fatal(err)
	}

	if err := getProductController.Handle("CHAIR-001"); err != nil {
		log.Fatal(err)
	}

	if err := listProductsController.Handle("Standard", true); err != nil {
		log.Fatal(err)
	}

	if err := getOrderController.Handle(convertPresenter.ViewModel().OrderID); err != nil {
		log.Fatal(err)
	}

	if err := listOrdersController.Handle(entities.OrderStatusShipped); err != nil {
		log.Fatal(err)
	}

	if err := getShipmentController.Handle(shipmentPresenter.ViewModel().ShipmentID); err != nil {
		log.Fatal(err)
	}

	if err := listShipmentsController.Handle(convertPresenter.ViewModel().OrderID); err != nil {
		log.Fatal(err)
	}

	if err := requestReturnController.Handle(convertPresenter.ViewModel().OrderID, "damaged item", "customer-001"); err != nil {
		log.Fatal(err)
	}

	if err := getReturnController.Handle(requestReturnPresenter.ViewModel().ReturnRequestID); err != nil {
		log.Fatal(err)
	}

	if err := listReturnController.Handle(entities.ReturnRequestStatusRequested); err != nil {
		log.Fatal(err)
	}

	if err := quoteConversionReportController.Handle(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(createPresenter.ViewModel().Message)
	fmt.Println(getCustomerPresenter.ViewModel().Message)
	fmt.Println(listCustomersPresenter.ViewModel().Message)
	fmt.Println(quoteConversionReportPresenter.ViewModel().Message)
	fmt.Println(addLinePresenter.ViewModel().Message)
	fmt.Println(submitPresenter.ViewModel().Message)
	fmt.Println(convertPresenter.ViewModel().Message)
	fmt.Println(capturePresenter.ViewModel().Message)
	fmt.Println(shipmentPresenter.ViewModel().Message)
	fmt.Println(getPresenter.ViewModel().Message)
	fmt.Println(listQuotesPresenter.ViewModel().Message)
	fmt.Println(getProductPresenter.ViewModel().Message)
	fmt.Println(listProductsPresenter.ViewModel().Message)
	fmt.Println(getOrderPresenter.ViewModel().Message)
	fmt.Println(listOrdersPresenter.ViewModel().Message)
	fmt.Println(getShipmentPresenter.ViewModel().Message)
	fmt.Println(listShipmentsPresenter.ViewModel().Message)
	fmt.Println(requestReturnPresenter.ViewModel().Message)
	fmt.Println(getReturnPresenter.ViewModel().Message)
	fmt.Println(listReturnPresenter.ViewModel().Message)
}
