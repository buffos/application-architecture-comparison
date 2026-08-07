package records

import (
	"errors"
	"sort"
	"strconv"
)

const (
	PluginTypePricing        = "pricing"
	PluginDiscountPercentKey = "discountPercent"
)

var (
	ErrPluginKeyRequired          = errors.New("plugin key is required")
	ErrPluginTypeRequired         = errors.New("plugin type is required")
	ErrPluginNotFound             = errors.New("plugin not found")
	ErrPluginConfigurationInvalid = errors.New("pricing plugin configuration is invalid")
)

// PluginRegistration is a persistence-aware extension registration record.
type PluginRegistration struct {
	db *Database

	Key     string
	Type    string
	Version string
	Enabled bool
	Config  map[string]string
}

// RegisterPlugin creates and persists a plugin registration with a defensive
// configuration copy.
func RegisterPlugin(db *Database, key string, pluginType string, version string, config map[string]string) (*PluginRegistration, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if key == "" {
		return nil, ErrPluginKeyRequired
	}
	if pluginType == "" {
		return nil, ErrPluginTypeRequired
	}

	plugin := &PluginRegistration{
		db:      db,
		Key:     key,
		Type:    pluginType,
		Version: version,
		Config:  clonePluginConfig(config),
	}
	if err := plugin.Save(); err != nil {
		return nil, err
	}
	return plugin, nil
}

// FindPlugin loads a plugin registration Active Record by key.
func FindPlugin(db *Database, key string) (*PluginRegistration, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	if key == "" {
		return nil, ErrPluginKeyRequired
	}
	row, ok := db.plugins[key]
	if !ok {
		return nil, ErrPluginNotFound
	}
	return &PluginRegistration{
		db:      db,
		Key:     row.Key,
		Type:    row.Type,
		Version: row.Version,
		Enabled: row.Enabled,
		Config:  clonePluginConfig(row.Config),
	}, nil
}

// Enable changes a plugin registration to enabled and persists it.
func (plugin *PluginRegistration) Enable() error {
	if plugin == nil || plugin.db == nil {
		return ErrDatabaseRequired
	}
	plugin.Enabled = true
	return plugin.Save()
}

// Disable changes a plugin registration to disabled and persists it.
func (plugin *PluginRegistration) Disable() error {
	if plugin == nil || plugin.db == nil {
		return ErrDatabaseRequired
	}
	plugin.Enabled = false
	return plugin.Save()
}

// Save writes the current plugin registration to the plugins table.
func (plugin *PluginRegistration) Save() error {
	if plugin == nil || plugin.db == nil {
		return ErrDatabaseRequired
	}
	if plugin.Key == "" {
		return ErrPluginKeyRequired
	}
	if plugin.Type == "" {
		return ErrPluginTypeRequired
	}
	plugin.db.plugins[plugin.Key] = pluginRow{
		Key:     plugin.Key,
		Type:    plugin.Type,
		Version: plugin.Version,
		Enabled: plugin.Enabled,
		Config:  clonePluginConfig(plugin.Config),
	}
	return nil
}

// PriceQuoteLine applies enabled pricing plugin contributions in key order
// and returns the base unit price, total discount, and line total.
func PriceQuoteLine(db *Database, product *Product, quantity int) (int, int, int, error) {
	if db == nil {
		return 0, 0, 0, ErrDatabaseRequired
	}
	if product == nil {
		return 0, 0, 0, ErrProductRequired
	}
	if quantity <= 0 {
		return 0, 0, 0, ErrQuantityInvalid
	}

	keys := make([]string, 0, len(db.plugins))
	for key, plugin := range db.plugins {
		if plugin.Enabled && plugin.Type == PluginTypePricing {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	discountPercent := 0
	for _, key := range keys {
		plugin := db.plugins[key]
		percent, err := strconv.Atoi(plugin.Config[PluginDiscountPercentKey])
		if err != nil || percent < 0 {
			return 0, 0, 0, ErrPluginConfigurationInvalid
		}
		discountPercent += percent
	}
	if discountPercent > 100 {
		discountPercent = 100
	}

	baseTotal := product.UnitPrice * quantity
	discountAmount := baseTotal * discountPercent / 100
	return product.UnitPrice, discountAmount, baseTotal - discountAmount, nil
}

func clonePluginConfig(config map[string]string) map[string]string {
	if config == nil {
		return nil
	}
	clone := make(map[string]string, len(config))
	for key, value := range config {
		clone[key] = value
	}
	return clone
}
