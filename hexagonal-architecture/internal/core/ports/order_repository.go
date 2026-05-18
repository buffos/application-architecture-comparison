package ports

import "hexagonal-architecture/internal/core/domain"

type OrderRepository interface {
	Save(order domain.Order) error
	FindByID(id string) (domain.Order, error)
}
