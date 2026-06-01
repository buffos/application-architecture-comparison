package application

type QuoteConversionReport struct {
	TotalQuotes     int
	ApprovedQuotes  int
	ConvertedQuotes int
	ConversionRate  float64
}

type QuoteConversionReportService struct {
	quotes QuoteFinder
	orders OrderFinder
}

func NewQuoteConversionReportService(quotes QuoteFinder, orders OrderFinder) QuoteConversionReportService {
	return QuoteConversionReportService{
		quotes: quotes,
		orders: orders,
	}
}

func (s QuoteConversionReportService) Execute() (QuoteConversionReport, error) {
	drafts, err := s.quotes.ListByStatus("")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	approved, err := s.quotes.ListByStatus("Approved")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	pendingApproval, err := s.quotes.ListByStatus("PendingApproval")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	pendingPayment, err := s.orders.ListByStatus("PendingPayment")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	paid, err := s.orders.ListByStatus("Paid")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	paymentReview, err := s.orders.ListByStatus("PaymentReview")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	shipped, err := s.orders.ListByStatus("Shipped")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	partiallyShipped, err := s.orders.ListByStatus("PartiallyShipped")
	if err != nil {
		return QuoteConversionReport{}, err
	}

	totalQuotes := len(drafts)
	approvedQuotes := len(approved) + len(pendingApproval)
	convertedQuotes := len(pendingPayment) + len(paymentReview) + len(paid) + len(partiallyShipped) + len(shipped)

	rate := 0.0
	if approvedQuotes > 0 {
		rate = float64(convertedQuotes) / float64(approvedQuotes)
	}

	return QuoteConversionReport{
		TotalQuotes:     totalQuotes,
		ApprovedQuotes:  approvedQuotes,
		ConvertedQuotes: convertedQuotes,
		ConversionRate:  rate,
	}, nil
}
