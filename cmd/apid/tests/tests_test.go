package tests

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/timurguseynov/go-wallet-api/cmd/apid/handlers"
	"github.com/timurguseynov/go-wallet-api/internal/platform/db"
	"github.com/timurguseynov/go-wallet-api/internal/platform/tests"
	"github.com/timurguseynov/go-wallet-api/internal/platform/user"

	"github.com/timurguseynov/go-wallet-api/internal/platform/web"
)

var (
	a                 *web.App
	test              *tests.Test
	createdUsersCount int
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
		addTestUser(ctx, dbConn, "Alex", i*100)
	}

	return nil
}

func addTestUser(ctx context.Context, dbConn *db.DB, name string, depositAmount int) error {
	u := user.User{
		Name: name + strconv.Itoa(depositAmount),
	}
	id, err := user.Insert(ctx, dbConn, u)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = user.DepositByID(ctx, dbConn, id, int64(depositAmount))
	if err != nil {
		return errors.Wrap(err, "")
	}

	createdUsersCount++

	return nil
}
