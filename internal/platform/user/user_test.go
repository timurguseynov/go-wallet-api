package user_test

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/timurguseynov/go-wallet-api/internal/platform/tests"
	"github.com/timurguseynov/go-wallet-api/internal/platform/user"
)

var test *tests.Test

// TestMain is the entry point for testing.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	test = tests.New()
	defer test.TearDown()
	return m.Run()
}

var (
	ctx            context.Context
	userID         string
	depositAmount  int64 = 10000
	withdrawAmount int64 = 5000
)

func TestUser(t *testing.T) {
	defer tests.Recover(t)
	ctx = tests.Context()

	t.Run("userInsert", userInsert)
	t.Run("userDepositByID", userDepositByID)
	t.Run("userWithdrawByID", userWithdrawByID)
	t.Run("userList", userList)
}

func userInsert(t *testing.T) {
	var err error
	u := user.User{
		Name: "Alex",
	}
	userID, err = user.Insert(ctx, test.MasterDB, u)
	require.NoError(t, err)
	require.NotEmpty(t, userID, "should have id generated")
}

func userDepositByID(t *testing.T) {
	err := user.DepositByID(ctx, test.MasterDB, userID, depositAmount)
	require.NoError(t, err)

	balance, err := user.GetBalanceByID(ctx, test.MasterDB, userID)
	require.NoError(t, err)
	require.Equal(t, depositAmount, balance)
}

func userWithdrawByID(t *testing.T) {
	err := user.WithdrawByID(ctx, test.MasterDB, userID, withdrawAmount)
	require.NoError(t, err)

	balance, err := user.GetBalanceByID(ctx, test.MasterDB, userID)
	require.NoError(t, err)
	require.Equal(t, depositAmount-withdrawAmount, balance)
}

func userList(t *testing.T) {
	users, err := user.List(ctx, test.MasterDB)
	require.NoError(t, err)
	require.True(t, len(users) > 2)
}

func userListLeaders(t *testing.T) {
	users, err := user.List(ctx, test.MasterDB)
	require.NoError(t, err)
	require.Equal(t, 1, len(users))

	require.True(t, sort.SliceIsSorted(users, func(i, j int) bool {
		return users[i].Balance < users[j].Balance
	}), "should be sorted by Balance")
}
