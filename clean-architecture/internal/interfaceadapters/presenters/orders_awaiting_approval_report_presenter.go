package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ApprovalQueueItemViewModel struct {
	QuoteID     string
	CustomerID  string
	LineCount   int
	TotalAmount int
}

type OrdersAwaitingApprovalReportViewModel struct {
	Message string
	Count   int
	Items   []ApprovalQueueItemViewModel
}

type OrdersAwaitingApprovalReportPresenter struct {
	viewModel OrdersAwaitingApprovalReportViewModel
}

func NewOrdersAwaitingApprovalReportPresenter() *OrdersAwaitingApprovalReportPresenter {
	return &OrdersAwaitingApprovalReportPresenter{}
}

func (p *OrdersAwaitingApprovalReportPresenter) Present(output usecases.OrdersAwaitingApprovalReportOutput) error {
	items := make([]ApprovalQueueItemViewModel, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, ApprovalQueueItemViewModel{
			QuoteID:     item.QuoteID,
			CustomerID:  item.CustomerID,
			LineCount:   item.LineCount,
			TotalAmount: item.TotalAmount,
		})
	}

	p.viewModel = OrdersAwaitingApprovalReportViewModel{
		Message: fmt.Sprintf("orders awaiting approval report: count=%d", output.Count),
		Count:   output.Count,
		Items:   items,
	}

	return nil
}

func (p *OrdersAwaitingApprovalReportPresenter) ViewModel() OrdersAwaitingApprovalReportViewModel {
	return p.viewModel
}
