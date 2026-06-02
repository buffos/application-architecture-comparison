package kernel

type Host struct {
	plugins           map[string]struct{}
	customerDirectory CustomerDirectory
	quoteService      QuoteService
	quoteReader       QuoteReader
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

func (h *Host) CustomerDirectory() (CustomerDirectory, error) {
	if h.customerDirectory == nil {
		return nil, ErrCustomerDirectoryNotRegistered
	}

	return h.customerDirectory, nil
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
