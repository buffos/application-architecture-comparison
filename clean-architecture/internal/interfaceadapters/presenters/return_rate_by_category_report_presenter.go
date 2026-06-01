package presenters

import (
	"fmt"
	"strings"

	"clean-architecture/internal/usecases"
)

type ReturnRateByCategoryItemViewModel struct {
	Category         string
	ShippedQuantity  int
	ReturnedQuantity int
	ReturnRate       float64
}

type ReturnRateByCategoryReportViewModel struct {
	Message    string
	Categories []ReturnRateByCategoryItemViewModel
}

type ReturnRateByCategoryReportPresenter struct {
	viewModel ReturnRateByCategoryReportViewModel
}

func NewReturnRateByCategoryReportPresenter() *ReturnRateByCategoryReportPresenter {
	return &ReturnRateByCategoryReportPresenter{}
}

func (p *ReturnRateByCategoryReportPresenter) Present(output usecases.ReturnRateByCategoryReportOutput) error {
	items := make([]ReturnRateByCategoryItemViewModel, 0, len(output.Categories))
	summary := make([]string, 0, len(output.Categories))
	for _, category := range output.Categories {
		items = append(items, ReturnRateByCategoryItemViewModel{
			Category:         category.Category,
			ShippedQuantity:  category.ShippedQuantity,
			ReturnedQuantity: category.ReturnedQuantity,
			ReturnRate:       category.ReturnRate,
		})
		summary = append(summary, fmt.Sprintf("%s=%.2f", category.Category, category.ReturnRate))
	}

	p.viewModel = ReturnRateByCategoryReportViewModel{
		Message:    fmt.Sprintf("return rate by category report: %s", strings.Join(summary, ", ")),
		Categories: items,
	}

	return nil
}

func (p *ReturnRateByCategoryReportPresenter) ViewModel() ReturnRateByCategoryReportViewModel {
	return p.viewModel
}
