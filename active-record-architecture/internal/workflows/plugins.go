package workflows

import "active-record-architecture/internal/records"

// RegisterPlugin creates a persisted plugin registration Active Record.
func RegisterPlugin(db *records.Database, key string, pluginType string, version string, config map[string]string) (*records.PluginRegistration, error) {
	return records.RegisterPlugin(db, key, pluginType, version, config)
}

// EnablePlugin loads and enables a plugin registration.
func EnablePlugin(db *records.Database, key string) (*records.PluginRegistration, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}
	plugin, err := records.FindPlugin(db, key)
	if err != nil {
		return nil, err
	}
	if err := plugin.Enable(); err != nil {
		return nil, err
	}
	return plugin, nil
}

// DisablePlugin loads and disables a plugin registration.
func DisablePlugin(db *records.Database, key string) (*records.PluginRegistration, error) {
	if db == nil {
		return nil, records.ErrDatabaseRequired
	}
	plugin, err := records.FindPlugin(db, key)
	if err != nil {
		return nil, err
	}
	if err := plugin.Disable(); err != nil {
		return nil, err
	}
	return plugin, nil
}

// PriceQuoteLine exposes the configured pricing calculation as a read-only
// application operation.
func PriceQuoteLine(db *records.Database, productSKU string, quantity int) (int, int, int, error) {
	if db == nil {
		return 0, 0, 0, records.ErrDatabaseRequired
	}
	product, err := records.FindProduct(db, productSKU)
	if err != nil {
		return 0, 0, 0, err
	}
	return records.PriceQuoteLine(db, product, quantity)
}
