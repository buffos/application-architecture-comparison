package plugins

import "errors"

var (
	ErrPluginNotFound       = errors.New("plugin not found")
	ErrPluginAlreadyExists  = errors.New("plugin already exists")
	ErrPluginNotPricingType = errors.New("plugin is not a pricing plugin")
)

const PluginTypePricing = "pricing"

type PluginRegistration struct {
	ID      string
	Type    string
	Enabled bool
}

type Repository interface {
	Save(plugin PluginRegistration) error
	FindByID(id string) (PluginRegistration, error)
	List() ([]PluginRegistration, error)
}
