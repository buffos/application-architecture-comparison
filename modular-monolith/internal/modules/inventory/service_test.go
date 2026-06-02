package inventory

import "testing"

type stubRepository struct {
	reserved []ReservationItem
	released []ReleaseItem
	err      error
}

func (r *stubRepository) Save(record StockRecord) error {
	return nil
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
