package tests

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/timurguseynov/go-wallet-api/cmd/apid/handlers"
	"github.com/timurguseynov/go-wallet-api/internal/platform/rest"
	"github.com/timurguseynov/go-wallet-api/internal/platform/tests"
)

var (
	a    *rest.App
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

	a = handlers.API(test.MasterDB).(*rest.App)

	return m.Run()
}
