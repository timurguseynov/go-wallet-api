package tests

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/timurguseynov/user-wallet-api/cmd/apid/handlers"
	"github.com/timurguseynov/user-wallet-api/internal/platform/db"
	"github.com/timurguseynov/user-wallet-api/internal/platform/tests"
	"github.com/timurguseynov/user-wallet-api/internal/platform/user"

	"github.com/timurguseynov/user-wallet-api/internal/platform/web"
)

var (
	a    *web.App
	test *tests.Test
)

// TestMain is the entry point for testing.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// TestHandlers is the entry point for testing all handlers
func TestHandlers(t *testing.T) {
	defer tests.Recover(t)

	log.SetOutput(ioutil.Discard)

	t.Run("users", RunTestUser)
	t.Run("notifier", RunTestNotifier)
}

func testMain(m *testing.M) int {
	test = tests.New()
	defer test.TearDown()

	a = handlers.API(test.MasterDB).(*web.App)

	return m.Run()
}

func addTestUsers(ctx context.Context, dbConn *db.DB) error {
	for i := 0; i < 10; i++ {
		u := user.User{
			Name: "Alex" + strconv.Itoa(i),
		}
		id, err := user.Insert(ctx, dbConn, u)
		if err != nil {
			return errors.Wrap(err, "")
		}

		err = user.DepositByID(ctx, dbConn, id, int64(i*100))
		if err != nil {
			return errors.Wrap(err, "")

		}
	}

	return nil
}
