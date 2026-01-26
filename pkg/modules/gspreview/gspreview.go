package gspreview

import (
	"context"
	"sync"
	"time"

	"github.com/gotenberg/gotenberg/v8/pkg/gotenberg"
	"github.com/gotenberg/gotenberg/v8/pkg/modules/api"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/single_threaded"
)

func init() {
	gotenberg.MustRegisterModule(new(Module))
}

type Module struct {
	pdfiumInstance pdfium.Pdfium
	pdfiumMutex    sync.Mutex
}

func (m *Module) Descriptor() gotenberg.ModuleDescriptor {
	return gotenberg.ModuleDescriptor{
		ID:  "gspreview",
		New: func() gotenberg.Module { return new(Module) },
	}
}

func (m *Module) Provision(ctx *gotenberg.Context) error { return nil }
func (m *Module) Validate() error                        { return nil }
func (m *Module) Debug() map[string]interface{}          { return nil }

func (m *Module) Routes() ([]api.Route, error) {
	return []api.Route{
		gspreviewRoute(m),
	}, nil
}

func (m *Module) Start() error {
	var err error
	pool := single_threaded.Init(single_threaded.Config{})
	m.pdfiumInstance, err = pool.GetInstance(time.Second * 10)
	return err
}

func (m *Module) Stop(ctx context.Context) error { return nil }

func (m *Module) StartupMessage() string {
	return "Starting gspreview (pdfium-bindings)"
}

var (
	_ gotenberg.Module = (*Module)(nil)
	_ gotenberg.App    = (*Module)(nil)
)
