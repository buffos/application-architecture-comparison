package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestEnabledPricingPluginChangesQuoteLineTotal(t *testing.T) {
	store := data.NewStore()
	product := data.Product{SKU: "sku-001", UnitPrice: 10000}
	if _, err := RegisterPlugin(store, "seasonal-10", PluginTypePricing, "1.0", map[string]string{PluginDiscountPercentKey: "10"}); err != nil {
		t.Fatalf("RegisterPlugin() error = %v", err)
	}
	if _, err := EnablePlugin(store, "seasonal-10"); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}

	unitPrice, discount, total, err := PriceQuoteLine(store, product, 2)
	if err != nil {
		t.Fatalf("PriceQuoteLine() error = %v", err)
	}
	if unitPrice != 10000 || discount != 2000 || total != 18000 {
		t.Fatalf("pricing = unit=%d discount=%d total=%d, want 10000/2000/18000", unitPrice, discount, total)
	}
}

func TestDisabledPricingPluginHasNoEffectAndMultiplePluginsAreCapped(t *testing.T) {
	store := data.NewStore()
	if _, err := RegisterPlugin(store, "b", PluginTypePricing, "1", map[string]string{PluginDiscountPercentKey: "70"}); err != nil {
		t.Fatalf("RegisterPlugin b error = %v", err)
	}
	if _, err := RegisterPlugin(store, "a", PluginTypePricing, "1", map[string]string{PluginDiscountPercentKey: "70"}); err != nil {
		t.Fatalf("RegisterPlugin a error = %v", err)
	}
	if _, err := EnablePlugin(store, "a"); err != nil {
		t.Fatalf("EnablePlugin a error = %v", err)
	}

	_, _, total, err := PriceQuoteLine(store, data.Product{UnitPrice: 10000}, 1)
	if err != nil {
		t.Fatalf("PriceQuoteLine() one plugin error = %v", err)
	}
	if total != 3000 {
		t.Fatalf("one-plugin total = %d, want 3000", total)
	}
	if _, err := EnablePlugin(store, "b"); err != nil {
		t.Fatalf("EnablePlugin b error = %v", err)
	}
	_, _, total, err = PriceQuoteLine(store, data.Product{UnitPrice: 10000}, 1)
	if err != nil {
		t.Fatalf("PriceQuoteLine() capped error = %v", err)
	}
	if total != 0 {
		t.Fatalf("capped total = %d, want 0", total)
	}
}

func TestAddQuoteLinePersistsPluginDiscount(t *testing.T) {
	store := data.NewStore()
	store.Quotes["quote-001"] = data.Quote{ID: "quote-001", Status: data.QuoteStatusDraft}
	store.Products["sku-001"] = data.Product{SKU: "sku-001", Name: "Desk", Active: true, UnitPrice: 10000}
	if _, err := RegisterPlugin(store, "seasonal-10", PluginTypePricing, "1.0", map[string]string{PluginDiscountPercentKey: "10"}); err != nil {
		t.Fatalf("RegisterPlugin() error = %v", err)
	}
	if _, err := EnablePlugin(store, "seasonal-10"); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}

	quote, err := AddQuoteLine(store, "quote-001", "sku-001", 2)
	if err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	if quote.Lines[0].DiscountAmount != 2000 || quote.Lines[0].LineTotal != 18000 {
		t.Fatalf("line = %#v, want discount 2000 total 18000", quote.Lines[0])
	}
}
