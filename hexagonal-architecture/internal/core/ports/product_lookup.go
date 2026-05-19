package ports

import "hexagonal-architecture/internal/core/domain"

type ProductLookup interface {
	FindBySKU(sku string) (domain.Product, error)
	List(category string, availableOnly bool) ([]domain.Product, error)
}
