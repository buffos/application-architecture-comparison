package application

import "layered-architecture/internal/domain"

type OrderRepository interface {
	Save(order domain.Order) error
	FindByID(id string) (domain.Order, error)
}

type OrderService struct {
	orderRepo OrderRepository
	quoteRepo QuoteRepository
	stockRepo StockRecordRepository
}

func NewOrderService(orderRepo OrderRepository, quoteRepo QuoteRepository, stockRepo StockRecordRepository) OrderService {
	return OrderService{
		orderRepo: orderRepo,
		quoteRepo: quoteRepo,
		stockRepo: stockRepo,
	}
}

func (s OrderService) ConvertQuoteToOrder(quoteID string) (domain.Order, error) {
	quote, err := s.quoteRepo.FindByID(quoteID)
	if err != nil {
		return domain.Order{}, err
	}

	if quote.Status != domain.QuoteStatusApproved {
		return domain.Order{}, domain.ErrQuoteNotApproved
	}

	stocks := make([]domain.StockRecord, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		stock, err := s.stockRepo.FindBySKU(line.SKU)
		if err != nil {
			return domain.Order{}, err
		}

		if stock.Available() < line.Quantity {
			return domain.Order{}, domain.ErrInsufficientStock
		}

		stocks = append(stocks, stock)
	}

	for i, line := range quote.Lines {
		stock := stocks[i]
		if err := stock.Reserve(line.Quantity); err != nil {
			return domain.Order{}, err
		}

		if err := s.stockRepo.Save(stock); err != nil {
			return domain.Order{}, err
		}
	}

	order, err := domain.NewOrderFromQuote(quote)
	if err != nil {
		return domain.Order{}, err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (s OrderService) GetOrder(id string) (domain.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s OrderService) CancelOrder(id string) (domain.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return domain.Order{}, err
	}

	if err := order.Cancel(); err != nil {
		return domain.Order{}, err
	}

	for _, line := range order.Lines {
		stock, err := s.stockRepo.FindBySKU(line.SKU)
		if err != nil {
			return domain.Order{}, err
		}

		if err := stock.ReleaseReserved(line.Quantity); err != nil {
			return domain.Order{}, err
		}

		if err := s.stockRepo.Save(stock); err != nil {
			return domain.Order{}, err
		}
	}

	if err := s.orderRepo.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
