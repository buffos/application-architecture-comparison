package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestEnabledPricingPluginChangesQuoteLineTotal(t *testing.T) {
	db, quote := setupDraftQuote(t, true)
	if _, err := RegisterPlugin(db, "seasonal-10", records.PluginTypePricing, "1.0", map[string]string{records.PluginDiscountPercentKey: "10"}); err != nil {
		t.Fatalf("RegisterPlugin() error = %v", err)
	}
	if _, err := EnablePlugin(db, "seasonal-10"); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}

	got, err := AddQuoteLine(db, quote.ID, "sku-001", 2)
	if err != nil {
		t.Fatalf("AddQuoteLine() error = %v", err)
	}
	if got.Lines[0].UnitPrice != 15000 || got.Lines[0].DiscountAmount != 3000 || got.Lines[0].LineTotal != 27000 {
		t.Fatalf("priced quote line = %#v, want unit 15000 discount 3000 total 27000", got.Lines[0])
	}
}

func TestDisabledPricingPluginHasNoEffectAndMultiplePluginsAreCapped(t *testing.T) {
	db := records.NewDatabase()
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 10000)
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}
	if _, err := RegisterPlugin(db, "b", records.PluginTypePricing, "1", map[string]string{records.PluginDiscountPercentKey: "70"}); err != nil {
		t.Fatalf("RegisterPlugin b error = %v", err)
	}
	if _, err := RegisterPlugin(db, "a", records.PluginTypePricing, "1", map[string]string{records.PluginDiscountPercentKey: "70"}); err != nil {
		t.Fatalf("RegisterPlugin a error = %v", err)
	}

	unitPrice, discount, total, err := PriceQuoteLine(db, product.SKU, 1)
	if err != nil {
		t.Fatalf("disabled pricing error = %v", err)
	}
	if unitPrice != 10000 || discount != 0 || total != 10000 {
		t.Fatalf("disabled pricing = unit %d discount %d total %d", unitPrice, discount, total)
	}
	if _, err := EnablePlugin(db, "a"); err != nil {
		t.Fatalf("EnablePlugin a error = %v", err)
	}
	_, _, total, err = PriceQuoteLine(db, product.SKU, 1)
	if err != nil {
		t.Fatalf("one-plugin pricing error = %v", err)
	}
	if total != 3000 {
		t.Fatalf("one-plugin total = %d, want 3000", total)
	}
	if _, err := EnablePlugin(db, "b"); err != nil {
		t.Fatalf("EnablePlugin b error = %v", err)
	}
	_, _, total, err = PriceQuoteLine(db, product.SKU, 1)
	if err != nil {
		t.Fatalf("capped pricing error = %v", err)
	}
	if total != 0 {
		t.Fatalf("capped total = %d, want 0", total)
	}
}

func TestPluginConfigurationPersistsAndCanBeDisabled(t *testing.T) {
	db := records.NewDatabase()
	config := map[string]string{records.PluginDiscountPercentKey: "10"}
	if _, err := RegisterPlugin(db, "seasonal-10", records.PluginTypePricing, "1.0", config); err != nil {
		t.Fatalf("RegisterPlugin() error = %v", err)
	}
	config[records.PluginDiscountPercentKey] = "99"
	saved, err := records.FindPlugin(db, "seasonal-10")
	if err != nil {
		t.Fatalf("FindPlugin() error = %v", err)
	}
	if saved.Config[records.PluginDiscountPercentKey] != "10" || saved.Enabled {
		t.Fatalf("saved plugin = %#v", saved)
	}
	if _, err := EnablePlugin(db, saved.Key); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}
	if _, err := DisablePlugin(db, saved.Key); err != nil {
		t.Fatalf("DisablePlugin() error = %v", err)
	}
	final, err := records.FindPlugin(db, saved.Key)
	if err != nil {
		t.Fatalf("FindPlugin() final error = %v", err)
	}
	if final.Enabled {
		t.Fatalf("plugin after disable = %#v, want disabled", final)
	}
}

func TestPricingRejectsInvalidEnabledConfiguration(t *testing.T) {
	db := records.NewDatabase()
	product := records.NewProduct(db, "sku-001", "Desk", "Standard", true, 10000)
	if err := product.Save(); err != nil {
		t.Fatalf("Product.Save() error = %v", err)
	}
	if _, err := RegisterPlugin(db, "bad", records.PluginTypePricing, "1", map[string]string{records.PluginDiscountPercentKey: "not-a-number"}); err != nil {
		t.Fatalf("RegisterPlugin() error = %v", err)
	}
	if _, err := EnablePlugin(db, "bad"); err != nil {
		t.Fatalf("EnablePlugin() error = %v", err)
	}
	if _, _, _, err := PriceQuoteLine(db, product.SKU, 1); err != records.ErrPluginConfigurationInvalid {
		t.Fatalf("invalid configuration error = %v, want %v", err, records.ErrPluginConfigurationInvalid)
	}
}
