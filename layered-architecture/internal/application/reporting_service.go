package application

import "layered-architecture/internal/domain"

type LowStockItemView struct {
	SKU       string
	OnHand    int
	Reserved  int
	Available int
}

type QuoteApprovalView struct {
	QuoteID       string
	CustomerID    string
	LineCount     int
	CurrentStatus string
}

type QuoteConversionReportView struct {
	TotalQuotes     int
	ConvertedQuotes int
	ConversionRate  float64
}

type ReportingQueryService struct {
	quoteRepo QuoteRepository
	orderRepo OrderRepository
	stockRepo StockRecordRepository
}

func NewReportingQueryService(quoteRepo QuoteRepository, orderRepo OrderRepository, stockRepo StockRecordRepository) ReportingQueryService {
	return ReportingQueryService{
		quoteRepo: quoteRepo,
		orderRepo: orderRepo,
		stockRepo: stockRepo,
	}
}

func (s ReportingQueryService) GetOrdersAwaitingApproval() ([]QuoteApprovalView, error) {
	quotes, err := s.quoteRepo.List()
	if err != nil {
		return nil, err
	}

	views := make([]QuoteApprovalView, 0)
	for _, quote := range quotes {
		if quote.Status == domain.QuoteStatusPendingApproval {
			views = append(views, QuoteApprovalView{
				QuoteID:       quote.ID,
				CustomerID:    quote.CustomerID,
				LineCount:     len(quote.Lines),
				CurrentStatus: quote.Status,
			})
		}
	}

	return views, nil
}

func (s ReportingQueryService) GetLowStockItems(threshold int) ([]LowStockItemView, error) {
	stocks, err := s.stockRepo.List()
	if err != nil {
		return nil, err
	}

	views := make([]LowStockItemView, 0)
	for _, stock := range stocks {
		if stock.Available() <= threshold {
			views = append(views, LowStockItemView{
				SKU:       stock.SKU,
				OnHand:    stock.OnHand,
				Reserved:  stock.Reserved,
				Available: stock.Available(),
			})
		}
	}

	return views, nil
}

func (s ReportingQueryService) GetQuoteConversionReport() (QuoteConversionReportView, error) {
	quotes, err := s.quoteRepo.List()
	if err != nil {
		return QuoteConversionReportView{}, err
	}

	orders, err := s.orderRepo.List()
	if err != nil {
		return QuoteConversionReportView{}, err
	}

	totalQuotes := len(quotes)
	convertedQuotes := len(orders)
	rate := 0.0
	if totalQuotes > 0 {
		rate = float64(convertedQuotes) / float64(totalQuotes)
	}

	return QuoteConversionReportView{
		TotalQuotes:     totalQuotes,
		ConvertedQuotes: convertedQuotes,
		ConversionRate:  rate,
	}, nil
}
