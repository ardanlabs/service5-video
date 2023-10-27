package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"testing"

	"github.com/ardanlabs/service/app/services/sales-api/v1/handlers"
	"github.com/ardanlabs/service/business/core/user"
	"github.com/ardanlabs/service/business/data/dbtest"
	"github.com/ardanlabs/service/business/data/order"
	v1 "github.com/ardanlabs/service/business/web/v1"
)

type seedData struct {
	users []user.User
}

// WebTests holds methods for each subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type WebTests struct {
	app        http.Handler
	userToken  string
	adminToken string
}

// Test_Web is the entry point for testing CRUD base web apis.
func Test_Web(t *testing.T) {
	t.Parallel()

	test := dbtest.NewTest(t, c)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.Teardown()
	}()

	api := test.CoreAPIs

	shutdown := make(chan os.Signal, 1)
	tests := WebTests{
		app: v1.APIMux(v1.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.V1.Auth,
			DB:       test.DB,
		}, handlers.Routes{}),
		userToken:  test.TokenV1("user@example.com", "gophers"),
		adminToken: test.TokenV1("admin@example.com", "gophers"),
	}

	// -------------------------------------------------------------------------

	seed := func(ctx context.Context, api dbtest.CoreAPIs) (seedData, error) {
		usrs, err := api.User.Query(ctx, user.QueryFilter{}, order.By{Field: user.OrderByName, Direction: order.ASC}, 1, 2)
		if err != nil {
			return seedData{}, fmt.Errorf("seeding users : %w", err)
		}

		sd := seedData{
			users: usrs,
		}

		return sd, nil
	}

	t.Log("Seeding data ...")

	sd, err := seed(context.Background(), api)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	t.Run("get200", tests.get200(sd))
}

func (wt *WebTests) get200(sd seedData) func(t *testing.T) {
	return func(t *testing.T) {
	}
}
