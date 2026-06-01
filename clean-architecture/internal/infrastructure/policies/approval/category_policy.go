package approval

import "clean-architecture/internal/entities"

type CategoryPolicy struct{}

func NewCategoryPolicy() CategoryPolicy {
	return CategoryPolicy{}
}

func (p CategoryPolicy) RequiresApproval(quote entities.Quote) (bool, error) {
	for _, line := range quote.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true, nil
		}
	}

	return false, nil
}
