package main

import (
	"fmt"
	"log"

	"transaction-script-architecture/internal/data"
	"transaction-script-architecture/internal/scripts"
)

func main() {
	store := data.NewStore()
	store.Customers["customer-001"] = data.Customer{
		ID:     "customer-001",
		Active: true,
	}
	store.Products["sku-001"] = data.Product{
		SKU:                 "sku-001",
		Name:                "Desk",
		Category:            "Standard",
		Active:              true,
		UnitPrice:           15000,
		StockShortagePolicy: data.StockShortageRejectOrder,
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 10}
	store.Products["sku-002"] = data.Product{
		SKU:       "sku-002",
		Name:      "Custom Desk",
		Category:  "CustomBuild",
		Active:    true,
		UnitPrice: 45000,
	}

	quote, err := scripts.CreateDraftQuote(store, "customer-001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"created draft quote: id=%s customer=%s status=%s\n",
		quote.ID,
		quote.CustomerID,
		quote.Status,
	)

	quote, err = scripts.AddQuoteLine(store, quote.ID, "sku-001", 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"added quote line: quote=%s sku=%s quantity=%d total=%d\n",
		quote.ID,
		quote.Lines[0].SKU,
		quote.Lines[0].Quantity,
		quote.Lines[0].LineTotal,
	)

	quote, err = scripts.SubmitQuoteForApproval(store, quote.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"submitted quote for approval: quote=%s status=%s\n",
		quote.ID,
		quote.Status,
	)

	order, err := scripts.ConvertQuoteToOrder(store, quote.ID, "sales-1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"converted quote: quote=%s order=%s status=%s total=%d\n",
		quote.ID,
		order.ID,
		order.Status,
		order.Total,
	)

	order, err = scripts.CapturePayment(store, order.ID, scripts.PaymentOutcomeAccept)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"captured payment: order=%s status=%s payment=%s\n",
		order.ID,
		order.Status,
		order.PaymentID,
	)

	shipment, err := scripts.CreateShipment(store, order.ID, "warehouse-1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"created shipment: shipment=%s order=%s status=%s\n",
		shipment.ID,
		shipment.OrderID,
		shipment.Status,
	)

	pendingQuote, err := scripts.CreateDraftQuote(store, "customer-001")
	if err != nil {
		log.Fatal(err)
	}

	pendingQuote, err = scripts.AddQuoteLine(store, pendingQuote.ID, "sku-002", 1)
	if err != nil {
		log.Fatal(err)
	}

	pendingQuote, err = scripts.SubmitQuoteForApproval(store, pendingQuote.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"submitted custom quote: quote=%s status=%s\n",
		pendingQuote.ID,
		pendingQuote.Status,
	)

	pendingQuote, err = scripts.ApproveQuote(store, pendingQuote.ID, "manager-1", "Approved after review")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"approved quote: quote=%s status=%s reviewedBy=%s comment=%s\n",
		pendingQuote.ID,
		pendingQuote.Status,
		pendingQuote.ReviewedBy,
		pendingQuote.DecisionComment,
	)

	rejectedQuote, err := scripts.CreateDraftQuote(store, "customer-001")
	if err != nil {
		log.Fatal(err)
	}

	rejectedQuote, err = scripts.AddQuoteLine(store, rejectedQuote.ID, "sku-002", 1)
	if err != nil {
		log.Fatal(err)
	}

	rejectedQuote, err = scripts.SubmitQuoteForApproval(store, rejectedQuote.ID)
	if err != nil {
		log.Fatal(err)
	}

	rejectedQuote, err = scripts.RejectQuote(store, rejectedQuote.ID, "manager-2", "Customer credit limit exceeded")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"rejected quote: quote=%s status=%s reviewedBy=%s comment=%s\n",
		rejectedQuote.ID,
		rejectedQuote.Status,
		rejectedQuote.ReviewedBy,
		rejectedQuote.DecisionComment,
	)

}
