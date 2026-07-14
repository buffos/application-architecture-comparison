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
	inventoryRelease  InventoryRelease
	inventoryRestock  InventoryRestock
	payments          PaymentCapture
	refunds           PaymentRefund
	shipments         ShipmentCreation
	orderService      OrderService
	returnableOrders  ReturnableOrderProvider
	returnEligibility ReturnEligibilityPolicy
	returnService     ReturnService
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

func (h *Host) ExposeInventoryRelease(release InventoryRelease) {
	h.inventoryRelease = release
}

func (h *Host) ExposeInventoryRestock(restock InventoryRestock) {
	h.inventoryRestock = restock
}

func (h *Host) ExposePaymentCapture(payments PaymentCapture) {
	h.payments = payments
}

func (h *Host) ExposePaymentRefund(refunds PaymentRefund) {
	h.refunds = refunds
}

func (h *Host) ExposeShipmentCreation(shipments ShipmentCreation) {
	h.shipments = shipments
}

func (h *Host) ExposeProductCatalog(catalog ProductCatalog) {
	h.productCatalog = catalog
}

func (h *Host) ExposeApprovalPolicy(policy ApprovalPolicy) {
	h.approvalPolicy = policy
}

func (h *Host) ExposeOrderService(service OrderService) {
	h.orderService = service
}

func (h *Host) ExposeReturnableOrderProvider(provider ReturnableOrderProvider) {
	h.returnableOrders = provider
}

func (h *Host) ExposeReturnEligibilityPolicy(policy ReturnEligibilityPolicy) {
	h.returnEligibility = policy
}

func (h *Host) ExposeReturnService(service ReturnService) {
	h.returnService = service
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

func (h *Host) InventoryRelease() (InventoryRelease, error) {
	if h.inventoryRelease == nil {
		return nil, ErrInventoryReleaseNotRegistered
	}

	return h.inventoryRelease, nil
}

func (h *Host) InventoryRestock() (InventoryRestock, error) {
	if h.inventoryRestock == nil {
		return nil, ErrInventoryRestockNotRegistered
	}

	return h.inventoryRestock, nil
}

func (h *Host) PaymentCapture() (PaymentCapture, error) {
	if h.payments == nil {
		return nil, ErrPaymentCaptureNotRegistered
	}

	return h.payments, nil
}

func (h *Host) PaymentRefund() (PaymentRefund, error) {
	if h.refunds == nil {
		return nil, ErrPaymentRefundNotRegistered
	}

	return h.refunds, nil
}

func (h *Host) ShipmentCreation() (ShipmentCreation, error) {
	if h.shipments == nil {
		return nil, ErrShipmentCreationNotRegistered
	}

	return h.shipments, nil
}

func (h *Host) OrderService() (OrderService, error) {
	if h.orderService == nil {
		return nil, ErrOrderServiceNotRegistered
	}

	return h.orderService, nil
}

func (h *Host) ReturnableOrderProvider() (ReturnableOrderProvider, error) {
	if h.returnableOrders == nil {
		return nil, ErrReturnableOrderProviderNotRegistered
	}

	return h.returnableOrders, nil
}

func (h *Host) ReturnEligibilityPolicy() (ReturnEligibilityPolicy, error) {
	if h.returnEligibility == nil {
		return nil, ErrReturnEligibilityPolicyNotRegistered
	}

	return h.returnEligibility, nil
}

func (h *Host) ReturnService() (ReturnService, error) {
	if h.returnService == nil {
		return nil, ErrReturnServiceNotRegistered
	}

	return h.returnService, nil
}
