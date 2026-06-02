package kernel

import "errors"

var ErrPluginAlreadyRegistered = errors.New("plugin already registered")
var ErrCustomerDirectoryNotRegistered = errors.New("customer directory capability not registered")
var ErrQuoteServiceNotRegistered = errors.New("quote service capability not registered")

type Plugin interface {
	ID() string
	Register(host *Host) error
}

type CustomerDirectory interface {
	RequireActiveCustomer(id string) error
}

type CreateDraftQuoteCommand struct {
	CustomerID string
}

type CreateDraftQuoteResult struct {
	QuoteID    string
	CustomerID string
	Status     string
}

type QuoteService interface {
	CreateDraftQuote(command CreateDraftQuoteCommand) (CreateDraftQuoteResult, error)
}
