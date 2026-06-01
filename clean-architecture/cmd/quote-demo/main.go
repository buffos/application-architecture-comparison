package main

import (
	"fmt"
	"log"

	"clean-architecture/internal/entities"
	"clean-architecture/internal/infrastructure/memory"
	approvalpolicy "clean-architecture/internal/infrastructure/policies/approval"
	paymentservice "clean-architecture/internal/infrastructure/services/payment"
	"clean-architecture/internal/interfaceadapters/controllers"
	"clean-architecture/internal/interfaceadapters/presenters"
	"clean-architecture/internal/usecases"
)

func main() {
	customerGateway := memory.NewCustomerGateway()
	quoteGateway := memory.NewQuoteGateway()
	orderGateway := memory.NewOrderGateway()
	shipmentGateway := memory.NewShipmentGateway()
	productGateway := memory.NewProductGateway()
	inventoryReservation := memory.NewInventoryReservation(map[string]int{
		"CHAIR-001": 5,
	})
	approvalPolicy := approvalpolicy.NewCategoryPolicy()
	paymentGateway := paymentservice.NewAcceptAllGateway()
	createPresenter := presenters.NewCreateDraftQuotePresenter()
	addLinePresenter := presenters.NewAddQuoteLinePresenter()
	submitPresenter := presenters.NewSubmitQuotePresenter()
	convertPresenter := presenters.NewConvertQuoteToOrderPresenter()
	capturePresenter := presenters.NewCapturePaymentPresenter()
	shipmentPresenter := presenters.NewCreateShipmentPresenter()
	getPresenter := presenters.NewGetQuotePresenter()

	if err := customerGateway.Save(entities.Customer{
		ID:     "customer-001",
		Active: true,
	}); err != nil {
		log.Fatal(err)
	}
	if err := productGateway.Save(entities.Product{
		SKU:       "CHAIR-001",
		Name:      "Office Chair",
		Category:  "Standard",
		BasePrice: 10000,
		Available: true,
	}); err != nil {
		log.Fatal(err)
	}

	createInteractor := usecases.NewCreateDraftQuoteInteractor(quoteGateway, customerGateway, createPresenter)
	createController := controllers.NewCreateDraftQuoteController(createInteractor)

	if err := createController.Handle("customer-001"); err != nil {
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

	convertInteractor := usecases.NewConvertQuoteToOrderInteractor(quoteGateway, orderGateway, inventoryReservation, convertPresenter)
	convertController := controllers.NewConvertQuoteToOrderController(convertInteractor)

	if err := convertController.Handle(createPresenter.ViewModel().QuoteID); err != nil {
		log.Fatal(err)
	}

	captureInteractor := usecases.NewCapturePaymentInteractor(orderGateway, paymentGateway, capturePresenter)
	captureController := controllers.NewCapturePaymentController(captureInteractor)

	if err := captureController.Handle(convertPresenter.ViewModel().OrderID); err != nil {
		log.Fatal(err)
	}

	createShipmentInteractor := usecases.NewCreateShipmentInteractor(orderGateway, shipmentGateway, shipmentPresenter)
	createShipmentController := controllers.NewCreateShipmentController(createShipmentInteractor)

	if err := createShipmentController.Handle(convertPresenter.ViewModel().OrderID); err != nil {
		log.Fatal(err)
	}

	getInteractor := usecases.NewGetQuoteInteractor(quoteGateway, getPresenter)
	getController := controllers.NewGetQuoteController(getInteractor)

	if err := getController.Handle(createPresenter.ViewModel().QuoteID); err != nil {
		log.Fatal(err)
	}

	fmt.Println(createPresenter.ViewModel().Message)
	fmt.Println(addLinePresenter.ViewModel().Message)
	fmt.Println(submitPresenter.ViewModel().Message)
	fmt.Println(convertPresenter.ViewModel().Message)
	fmt.Println(capturePresenter.ViewModel().Message)
	fmt.Println(shipmentPresenter.ViewModel().Message)
	fmt.Println(getPresenter.ViewModel().Message)
}
