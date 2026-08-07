package scripts

import (
	"testing"

	"transaction-script-architecture/internal/data"
)

func TestEvaluateReturnEligibilityRejectsClearance(t *testing.T) {
	decision := EvaluateReturnEligibility(
		data.Order{Lines: []data.OrderLine{{ID: "line-1", ProductCategory: "Clearance"}}},
		data.ReturnRequest{Lines: []data.ReturnLine{{OrderLineID: "line-1", ProductCategory: "Clearance"}}},
	)

	if decision.Eligible {
		t.Fatal("Eligible = true, want false")
	}
	if len(decision.Reasons) != 1 || decision.Reasons[0] != ReturnEligibilityReasonClearance {
		t.Fatalf("Reasons = %#v, want clearance reason", decision.Reasons)
	}
}

func TestEvaluateReturnEligibilityAllowsStandardProduct(t *testing.T) {
	decision := EvaluateReturnEligibility(
		data.Order{Lines: []data.OrderLine{{ID: "line-1", ProductCategory: "Standard"}}},
		data.ReturnRequest{Lines: []data.ReturnLine{{OrderLineID: "line-1", ProductCategory: "Standard"}}},
	)

	if !decision.Eligible || len(decision.Reasons) != 0 {
		t.Fatalf("decision = %#v, want eligible without reasons", decision)
	}
}

func TestAcceptReturnRejectsClearanceBeforeSideEffects(t *testing.T) {
	store := data.NewStore()
	store.Orders["order-001"] = data.Order{
		ID:     "order-001",
		Status: data.OrderStatusShipped,
		Lines: []data.OrderLine{{
			ID:              "order-001-line-001",
			SKU:             "sku-001",
			ProductCategory: "Clearance",
			ShippedQuantity: 1,
			UnitPrice:       10000,
		}},
	}
	store.Stocks["sku-001"] = data.StockRecord{SKU: "sku-001", OnHand: 0}
	store.Returns["return-001"] = data.ReturnRequest{
		ID:      "return-001",
		OrderID: "order-001",
		Status:  data.ReturnStatusRequested,
		Lines: []data.ReturnLine{{
			OrderLineID:     "order-001-line-001",
			SKU:             "sku-001",
			ProductCategory: "Clearance",
			Quantity:        1,
		}},
	}

	_, err := AcceptReturn(store, "return-001")
	if err != ErrReturnNotEligible {
		t.Fatalf("error = %v, want %v", err, ErrReturnNotEligible)
	}
	if store.Returns["return-001"].Status != data.ReturnStatusRequested {
		t.Fatalf("status = %q, want requested", store.Returns["return-001"].Status)
	}
	if len(store.Refunds) != 0 || store.Stocks["sku-001"].OnHand != 0 {
		t.Fatalf("side effects occurred: refunds=%d stock=%#v", len(store.Refunds), store.Stocks["sku-001"])
	}
}
