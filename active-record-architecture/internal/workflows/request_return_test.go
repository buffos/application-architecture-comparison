package workflows

import (
	"testing"

	"active-record-architecture/internal/records"
)

func TestRequestReturnCreatesRequestedReturnAndNotStartedRefund(t *testing.T) {
	db, order := readyToShip(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}

	request, err := RequestReturn(db, order.ID, nil, "damaged on arrival", "customer-1")
	if err != nil {
		t.Fatalf("RequestReturn() error = %v", err)
	}
	if request.ID != "return-001" || request.Status != records.ReturnStatusRequested || request.RefundStatus != records.RefundStatusNotStarted {
		t.Fatalf("return request = %#v", request)
	}
	if request.RefundID != "refund-001" || request.RefundAmount != order.Total {
		t.Fatalf("return refund fields = %#v", request)
	}
	if len(request.Lines) != 1 || request.Lines[0].Quantity != 1 || request.Lines[0].SKU != "sku-001" {
		t.Fatalf("return lines = %#v", request.Lines)
	}

	savedRequest, err := records.FindReturnRequest(db, request.ID)
	if err != nil {
		t.Fatalf("FindReturnRequest() error = %v", err)
	}
	if savedRequest.Lines[0] != request.Lines[0] {
		t.Fatalf("saved return lines = %#v, want %#v", savedRequest.Lines, request.Lines)
	}
	savedRefund, err := records.FindRefund(db, request.RefundID)
	if err != nil {
		t.Fatalf("FindRefund() error = %v", err)
	}
	if savedRefund.ReturnRequestID != request.ID || savedRefund.Amount != order.Total || savedRefund.Status != records.RefundStatusNotStarted {
		t.Fatalf("saved refund = %#v", savedRefund)
	}
	stock, err := records.FindStock(db, "sku-001")
	if err != nil {
		t.Fatalf("FindStock() error = %v", err)
	}
	if stock.OnHand != 4 || stock.Reserved != 0 {
		t.Fatalf("stock after request = %#v, want on-hand 4 reserved 0", stock)
	}
}

func TestRequestReturnAllowsExplicitQuantityAndRejectsExcess(t *testing.T) {
	db, order := readyToShip(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	lines := []records.ReturnLine{{OrderLineID: order.Lines[0].ID, Quantity: 1}}
	request, err := RequestReturn(db, order.ID, lines, "partial return", "customer-1")
	if err != nil {
		t.Fatalf("explicit RequestReturn() error = %v", err)
	}
	if request.Lines[0].Quantity != 1 || request.RefundAmount != order.Lines[0].UnitPrice {
		t.Fatalf("explicit return = %#v", request)
	}

	db, order = readyToShip(t)
	if _, err := CreateShipment(db, order.ID, "warehouse-1"); err != nil {
		t.Fatalf("CreateShipment() error = %v", err)
	}
	tooMany := []records.ReturnLine{{OrderLineID: order.Lines[0].ID, Quantity: 2}}
	if _, err := RequestReturn(db, order.ID, tooMany, "too many", "customer-1"); err != records.ErrReturnLinesInvalid {
		t.Fatalf("excess quantity error = %v, want %v", err, records.ErrReturnLinesInvalid)
	}
	if _, err := records.FindReturnRequest(db, "return-001"); err != records.ErrReturnNotFound {
		t.Fatalf("return after rejected request = %v, want %v", err, records.ErrReturnNotFound)
	}
}

func TestRequestReturnRequiresShippedOrder(t *testing.T) {
	db, order := readyOrder(t)
	if _, err := RequestReturn(db, order.ID, nil, "not shipped", "customer-1"); err != records.ErrOrderNotReturnable {
		t.Fatalf("not-shipped error = %v, want %v", err, records.ErrOrderNotReturnable)
	}
}
