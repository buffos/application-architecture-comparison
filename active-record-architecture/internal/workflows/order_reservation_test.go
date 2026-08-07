package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestConvertQuoteToOrderRejectsHardShortageWithoutPersistence(t *testing.T) {
	db, quote := quoteWithLine(t, "Standard")
	stock := records.NewStockRecord(db, "sku-001", 0, 0, 2)
	if err := stock.Save(); err != nil {
		t.Fatalf("StockRecord.Save() error = %v", err)
	}

	_, err := SubmitQuoteForApproval(db, quote.ID)
	if err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	_, err = ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != records.ErrInsufficientStock {
		t.Fatalf("error = %v, want %v", err, records.ErrInsufficientStock)
	}

	savedQuote, err := records.FindQuote(db, quote.ID)
	if err != nil {
		t.Fatalf("FindQuote() error = %v", err)
	}
	if savedQuote.Status != records.QuoteStatusApproved {
		t.Fatalf("quote status = %q, want approved", savedQuote.Status)
	}
	if _, err := records.FindOrder(db, "order-001"); err != records.ErrOrderNotFound {
		t.Fatalf("FindOrder() error = %v, want %v", err, records.ErrOrderNotFound)
	}
	savedStock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if savedStock.Reserved != 0 {
		t.Fatalf("reserved stock = %d, want 0", savedStock.Reserved)
	}
}

func TestConvertQuoteToOrderAllowsConfiguredBackorder(t *testing.T) {
	db, quote := quoteWithLine(t, "Standard")
	product, err := records.FindProduct(db, "sku-001")
	if err != nil {
		t.Fatalf("FindProduct() error = %v", err)
	}
	product.StockShortagePolicy = records.StockShortageAllowBackorder
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}
	stock := records.NewStockRecord(db, "sku-001", 0, 0, 2)
	if err := stock.Save(); err != nil {
		t.Fatalf("StockRecord.Save() error = %v", err)
	}

	if _, err := SubmitQuoteForApproval(db, quote.ID); err != nil {
		t.Fatalf("SubmitQuoteForApproval() error = %v", err)
	}
	order, err := ConvertQuoteToOrder(db, quote.ID, "sales-1")
	if err != nil {
		t.Fatalf("ConvertQuoteToOrder() error = %v", err)
	}
	if order.Status != records.OrderStatusBackordered {
		t.Fatalf("status = %q, want %q", order.Status, records.OrderStatusBackordered)
	}
	savedStock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if savedStock.Reserved != 0 {
		t.Fatalf("reserved stock = %d, want 0", savedStock.Reserved)
	}
}
