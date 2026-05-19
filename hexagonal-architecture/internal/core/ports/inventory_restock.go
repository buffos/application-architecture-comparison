package ports

import "hexagonal-architecture/internal/core/domain"

type InventoryRestock interface {
	Restock(lines []domain.ReservationLine) error
}
