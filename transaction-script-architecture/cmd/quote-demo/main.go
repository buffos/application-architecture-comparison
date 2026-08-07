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
		SKU:       "sku-001",
		Name:      "Desk",
		Category:  "Standard",
		Active:    true,
		UnitPrice: 15000,
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

}
