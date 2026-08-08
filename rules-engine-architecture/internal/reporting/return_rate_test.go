package reporting

import (
	"math"
	"testing"
)

func TestBuildReturnRateByCategoryGroupsAttempts(t *testing.T) {
	rows := BuildReturnRateByCategory([]ReturnRecord{
		{ProductCategory: "Standard", Accepted: true},
		{ProductCategory: "Standard", Accepted: false},
		{ProductCategory: "Clearance", Accepted: false},
		{ProductCategory: "Standard", Accepted: true},
	})

	if len(rows) != 2 || rows[0].ProductCategory != "Clearance" || rows[1].ProductCategory != "Standard" {
		t.Fatalf("expected sorted category rows, got %+v", rows)
	}
	if rows[1].AttemptedReturns != 3 || rows[1].AcceptedReturns != 2 {
		t.Fatalf("unexpected standard counts: %+v", rows[1])
	}
	if math.Abs(rows[1].ReturnRatePercent-66.6666667) > 0.0001 {
		t.Fatalf("expected standard rate 66.67%%, got %f", rows[1].ReturnRatePercent)
	}
}

func TestBuildReturnRateByCategoryHandlesEmptyCategory(t *testing.T) {
	rows := BuildReturnRateByCategory([]ReturnRecord{{Accepted: true}})

	if len(rows) != 1 || rows[0].ProductCategory != "Unknown" {
		t.Fatalf("expected unknown category row, got %+v", rows)
	}
}
