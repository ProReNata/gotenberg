package minibytes

import (
	"github.com/gotenberg/gotenberg/v8/pkg/gotenberg"
	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
)

func init() {
	gotenberg.MustRegisterModule(new(Module))
}

// TODO: add desired output dimensions as an option
type Module struct{}

func (m *Module) Descriptor() gotenberg.ModuleDescriptor {
	return gotenberg.ModuleDescriptor{
		ID:  "gspreview",
		New: func() gotenberg.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx *gotenberg.Context) error { return nil }
func (m *Module) Validate() error                        { return nil }
func (m *Module) Debug() map[string]interface{}          { return nil }

// Routes returns the HTTP routes.
func (mod *Module) Routes() ([]api.Route, error) {
	return []api.Route{
		gspreviewRoute(),
	}, nil
}
