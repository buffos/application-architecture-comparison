package approval

import "onion-architecture/internal/domain"

type CategoryPolicy struct{}

func NewCategoryPolicy() CategoryPolicy {
	return CategoryPolicy{}
}

func (p CategoryPolicy) RequiresApproval(quote domain.Quote) (bool, error) {
	for _, line := range quote.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true, nil
		}
	}

	return false, nil
}
