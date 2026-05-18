package ports

import "hexagonal-architecture/internal/core/domain"

type InventoryReservation interface {
	Reserve(lines []domain.ReservationLine) error
}
