package router_test

import (
	"testing"

	"github.com/elopez00/scale-backend/cmd/api/router"
	"github.com/elopez00/scale-backend/pkg/test"
)

func TestRouterGet(t *testing.T) {
	app := test.GetMockApp()
	if router.Get(app) == nil {
		t.Error("Router did not return a valid httprouter handle")
	}
}