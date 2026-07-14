package kernel

type Host struct {
	plugins           map[string]struct{}
	customerDirectory CustomerDirectory
	productCatalog    ProductCatalog
	approvalPolicy    ApprovalPolicy
	quoteService      QuoteService
	quoteReader       QuoteReader
	approvedQuotes    ApprovedQuoteProvider
	inventory         InventoryReservation
	payments          PaymentCapture
	orderService      OrderService
}

func NewHost() *Host {
	return &Host{
		plugins: make(map[string]struct{}),
	}
}

func (h *Host) RegisterPlugin(plugin Plugin) error {
	if _, exists := h.plugins[plugin.ID()]; exists {
		return ErrPluginAlreadyRegistered
	}

	if err := plugin.Register(h); err != nil {
		return err
	}

	h.plugins[plugin.ID()] = struct{}{}
	return nil
}

func (h *Host) ExposeCustomerDirectory(directory CustomerDirectory) {
	h.customerDirectory = directory
}

func (h *Host) ExposeQuoteService(service QuoteService) {
	h.quoteService = service
}

func (h *Host) ExposeQuoteReader(reader QuoteReader) {
	h.quoteReader = reader
}

func (h *Host) ExposeApprovedQuoteProvider(provider ApprovedQuoteProvider) {
	h.approvedQuotes = provider
}

func (h *Host) ExposeInventoryReservation(reservation InventoryReservation) {
	h.inventory = reservation
}

func (h *Host) ExposePaymentCapture(payments PaymentCapture) {
	h.payments = payments
}

func (h *Host) ExposeProductCatalog(catalog ProductCatalog) {
	h.productCatalog = catalog
}

func (h *Host) ExposeApprovalPolicy(policy ApprovalPolicy) {
	h.approvalPolicy = policy
}

func (h *Host) CustomerDirectory() (CustomerDirectory, error) {
	if h.customerDirectory == nil {
		return nil, ErrCustomerDirectoryNotRegistered
	}

	return h.customerDirectory, nil
}

func (h *Host) ProductCatalog() (ProductCatalog, error) {
	if h.productCatalog == nil {
		return nil, ErrProductCatalogNotRegistered
	}

	return h.productCatalog, nil
}

func (h *Host) ApprovalPolicy() (ApprovalPolicy, error) {
	if h.approvalPolicy == nil {
		return nil, ErrApprovalPolicyNotRegistered
	}

	return h.approvalPolicy, nil
}

func (h *Host) QuoteService() (QuoteService, error) {
	if h.quoteService == nil {
		return nil, ErrQuoteServiceNotRegistered
	}

	return h.quoteService, nil
}

func (h *Host) QuoteReader() (QuoteReader, error) {
	if h.quoteReader == nil {
		return nil, ErrQuoteReaderNotRegistered
	}

	return h.quoteReader, nil
}

func (h *Host) ApprovedQuoteProvider() (ApprovedQuoteProvider, error) {
	if h.approvedQuotes == nil {
		return nil, ErrApprovedQuoteProviderNotRegistered
	}

	return h.approvedQuotes, nil
}

func (h *Host) InventoryReservation() (InventoryReservation, error) {
	if h.inventory == nil {
		return nil, ErrInventoryReservationNotRegistered
	}

	return h.inventory, nil
}

func (h *Host) PaymentCapture() (PaymentCapture, error) {
	if h.payments == nil {
		return nil, ErrPaymentCaptureNotRegistered
	}

	return h.payments, nil
}

func (h *Host) ExposeOrderService(service OrderService) {
	h.orderService = service
}

func (h *Host) OrderService() (OrderService, error) {
	if h.orderService == nil {
		return nil, ErrOrderServiceNotRegistered
	}

	return h.orderService, nil
}
