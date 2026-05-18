package ports

import "hexagonal-architecture/internal/core/domain"

type InventoryRelease interface {
	Release(lines []domain.ReservationLine) error
}
