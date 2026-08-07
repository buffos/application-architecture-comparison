package data

// PluginRegistration is a passive extension registration record.
type PluginRegistration struct {
	Key     string
	Type    string
	Version string
	Enabled bool
	Config  map[string]string
}
