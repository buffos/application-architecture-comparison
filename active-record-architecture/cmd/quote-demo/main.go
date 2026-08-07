package main

import (
	"fmt"
	"log"

	"active-record-architecture/internal/records"
	"active-record-architecture/internal/workflows"
)

func main() {
	db := records.NewDatabase()

	customer := records.NewCustomer(db, "customer-001", true)
	if err := customer.Save(); err != nil {
		log.Fatal(err)
	}

	loadedCustomer, err := records.FindCustomer(db, customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 15000)
	if err := product.Save(); err != nil {
		log.Fatal(err)
	}
	stock := records.NewStockRecord(db, product.SKU, 10, 0, 2)
	if err := stock.Save(); err != nil {
		log.Fatal(err)
	}

	quote, err := records.NewDraftQuote(db, loadedCustomer.ID)
	if err != nil {
		log.Fatal(err)
	}
	if err := quote.Save(); err != nil {
		log.Fatal(err)
	}

	quote, err = workflows.AddQuoteLine(db, quote.ID, product.SKU, 2)
	if err != nil {
		log.Fatal(err)
	}

	quote, err = workflows.SubmitQuoteForApproval(db, quote.ID)
	if err != nil {
		log.Fatal(err)
	}

	reloadedQuote, err := records.FindQuote(db, quote.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"created draft quote: id=%s customer=%s status=%s\n",
		reloadedQuote.ID,
		reloadedQuote.CustomerID,
		reloadedQuote.Status,
	)
	fmt.Printf(
		"added quote line: quote=%s sku=%s quantity=%d total=%d\n",
		reloadedQuote.ID,
		reloadedQuote.Lines[0].SKU,
		reloadedQuote.Lines[0].Quantity,
		reloadedQuote.Lines[0].LineTotal,
	)
	fmt.Printf(
		"submitted quote for approval: quote=%s status=%s\n",
		reloadedQuote.ID,
		reloadedQuote.Status,
	)

	order, err := workflows.ConvertQuoteToOrder(db, reloadedQuote.ID, "sales-1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"converted quote: quote=%s order=%s status=%s total=%d\n",
		reloadedQuote.ID,
		order.ID,
		order.Status,
		order.Total,
	)

	customProduct := records.NewProduct(db, "sku-002", "Custom Desk", "CustomBuild", true, 45000)
	if err := customProduct.Save(); err != nil {
		log.Fatal(err)
	}
	pendingQuote, err := records.NewDraftQuote(db, loadedCustomer.ID)
	if err != nil {
		log.Fatal(err)
	}
	if err := pendingQuote.Save(); err != nil {
		log.Fatal(err)
	}
	if _, err := workflows.AddQuoteLine(db, pendingQuote.ID, customProduct.SKU, 1); err != nil {
		log.Fatal(err)
	}
	if _, err := workflows.SubmitQuoteForApproval(db, pendingQuote.ID); err != nil {
		log.Fatal(err)
	}
	rejectedQuote, err := workflows.RejectQuote(db, pendingQuote.ID, "manager-2", "Customer credit limit exceeded")
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
