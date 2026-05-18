package ports

import "hexagonal-architecture/internal/core/domain"

type InventoryConsumption interface {
	Consume(lines []domain.ReservationLine) error
}
