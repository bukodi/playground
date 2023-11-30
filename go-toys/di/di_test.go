package di

import (
	"github.com/samber/do"
	"testing"
)

type dbSvc struct {
	url string
}

type httpSvc struct {
	basePath string
	dbSvc    *dbSvc
}

func TestDIEnv(t *testing.T) {
	srv := do.New()

	do.Provide(srv, func(srv *do.Injector) (*dbSvc, error) {
		return &dbSvc{}, nil
	})

	do.Provide(srv, func(srv *do.Injector) (*httpSvc, error) {
		dbSvc, err := do.Invoke[*dbSvc](srv)
		if err != nil {
			return nil, err
		}
		return &httpSvc{
			basePath: "/api",
			dbSvc:    dbSvc,
		}, nil
	})

	for _, svcName := range srv.ListProvidedServices() {
		t.Logf("provided service: %s", svcName)
	}
	for _, svcName := range srv.ListInvokedServices() {
		t.Logf("invoked service: %s", svcName)
	}
}
