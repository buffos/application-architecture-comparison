package approval

import "hexagonal-architecture/internal/core/domain"

type CategoryApprovalPolicy struct{}

func NewCategoryApprovalPolicy() CategoryApprovalPolicy {
	return CategoryApprovalPolicy{}
}

func (CategoryApprovalPolicy) RequiresApproval(quote domain.Quote) (bool, error) {
	for _, line := range quote.Lines {
		if line.ProductCategory == "CustomBuild" {
			return true, nil
		}
	}

	return false, nil
}
