package kernel

import "testing"

type stubPlugin struct {
	id         string
	registerFn func(*Host) error
}

func (p stubPlugin) ID() string {
	return p.id
}

func (p stubPlugin) Register(host *Host) error {
	return p.registerFn(host)
}

func TestRegisterPluginRejectsDuplicates(t *testing.T) {
	host := NewHost()
	plugin := stubPlugin{
		id: "customers",
		registerFn: func(host *Host) error {
			return nil
		},
	}

	if err := host.RegisterPlugin(plugin); err != nil {
		t.Fatalf("expected first registration to succeed, got %v", err)
	}

	if err := host.RegisterPlugin(plugin); err != ErrPluginAlreadyRegistered {
		t.Fatalf("expected duplicate registration error, got %v", err)
	}
}
