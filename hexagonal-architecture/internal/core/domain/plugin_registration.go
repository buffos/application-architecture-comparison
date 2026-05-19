package domain

import "errors"

const PluginTypePricing = "Pricing"

const PluginStatusDisabled = "Disabled"
const PluginStatusEnabled = "Enabled"

var ErrPluginNotFound = errors.New("plugin not found")
var ErrPluginAlreadyEnabled = errors.New("plugin already enabled")
var ErrPluginKeyRequired = errors.New("plugin key is required")
var ErrPluginDiscountInvalid = errors.New("plugin discount percent must be between 1 and 99")

type PluginRegistration struct {
	Key             string
	Type            string
	Status          string
	DiscountPercent int
}

func NewPricingPlugin(key string, discountPercent int) (PluginRegistration, error) {
	if key == "" {
		return PluginRegistration{}, ErrPluginKeyRequired
	}
	if discountPercent <= 0 || discountPercent >= 100 {
		return PluginRegistration{}, ErrPluginDiscountInvalid
	}

	return PluginRegistration{
		Key:             key,
		Type:            PluginTypePricing,
		Status:          PluginStatusDisabled,
		DiscountPercent: discountPercent,
	}, nil
}

func (p *PluginRegistration) Enable() error {
	if p.Status == PluginStatusEnabled {
		return ErrPluginAlreadyEnabled
	}

	p.Status = PluginStatusEnabled
	return nil
}
