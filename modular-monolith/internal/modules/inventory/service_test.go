package inventory

import "testing"

type stubRepository struct {
	records   []StockRecord
	reserved  []ReservationItem
	released  []ReleaseItem
	restocked []RestockItem
	err       error
}

func (r *stubRepository) Save(record StockRecord) error {
	return nil
}

func (r *stubRepository) List() ([]StockRecord, error) {
	if r.err != nil {
		return nil, r.err
	}

	return append([]StockRecord(nil), r.records...), nil
}

func (r *stubRepository) Reserve(items []ReservationItem) error {
	if r.err != nil {
		return r.err
	}

	r.reserved = append([]ReservationItem(nil), items...)
	return nil
}

func (r *stubRepository) Release(items []ReleaseItem) error {
	if r.err != nil {
		return r.err
	}

	r.released = append([]ReleaseItem(nil), items...)
	return nil
}

func (r *stubRepository) Restock(items []RestockItem) error {
	if r.err != nil {
		return r.err
	}

	r.restocked = append([]RestockItem(nil), items...)
	return nil
}

func TestReservePassesReservationItemsToRepository(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	err := service.Reserve([]ReservationItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repository.reserved) != 1 {
		t.Fatalf("expected 1 reserved item, got %d", len(repository.reserved))
	}
}

func TestReserveRejectsNonPositiveQuantity(t *testing.T) {
	service := NewService(&stubRepository{})

	err := service.Reserve([]ReservationItem{
		{ProductSKU: "sku-001", Quantity: 0},
	})
	if err != ErrReservationQuantityMustBePositive {
		t.Fatalf("expected %v, got %v", ErrReservationQuantityMustBePositive, err)
	}
}

func TestReleasePassesReleaseItemsToRepository(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	err := service.Release([]ReleaseItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repository.released) != 1 {
		t.Fatalf("expected 1 released item, got %d", len(repository.released))
	}
}

func TestRestockPassesRestockItemsToRepository(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	err := service.Restock([]RestockItem{
		{ProductSKU: "sku-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repository.restocked) != 1 {
		t.Fatalf("expected 1 restocked item, got %d", len(repository.restocked))
	}
}

func TestListStockReturnsSnapshots(t *testing.T) {
	repository := &stubRepository{
		records: []StockRecord{
			{ProductSKU: "sku-001", Available: 3},
			{ProductSKU: "sku-002", Available: 9},
		},
	}
	service := NewService(repository)

	stock, err := service.ListStock()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stock) != 2 || stock[0].ProductSKU != "sku-001" {
		t.Fatalf("expected stock snapshots, got %+v", stock)
	}
}
