package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type LowStockItemViewModel struct {
	SKU      string
	Quantity int
}

type LowStockItemsReportViewModel struct {
	Message   string
	Threshold int
	Count     int
	Items     []LowStockItemViewModel
}

type LowStockItemsReportPresenter struct {
	viewModel LowStockItemsReportViewModel
}

func NewLowStockItemsReportPresenter() *LowStockItemsReportPresenter {
	return &LowStockItemsReportPresenter{}
}

func (p *LowStockItemsReportPresenter) Present(output usecases.LowStockItemsReportOutput) error {
	items := make([]LowStockItemViewModel, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, LowStockItemViewModel{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	p.viewModel = LowStockItemsReportViewModel{
		Message:   fmt.Sprintf("low stock items report: threshold=%d count=%d", output.Threshold, output.Count),
		Threshold: output.Threshold,
		Count:     output.Count,
		Items:     items,
	}

	return nil
}

func (p *LowStockItemsReportPresenter) ViewModel() LowStockItemsReportViewModel {
	return p.viewModel
}
