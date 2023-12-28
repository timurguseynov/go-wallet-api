package tests

import (
	"context"
	"log"
	"strconv"

	"github.com/pkg/errors"
	"github.com/timurguseynov/go-wallet-api/internal/db"
	"github.com/timurguseynov/go-wallet-api/internal/user"
)

func mustSeed(ctx context.Context, dbConn *db.DB) {
	err := seed10Users(ctx, dbConn)
	if err != nil {
		log.Fatal("couldn't seed users")
	}
}

func SeedUser(ctx context.Context, dbConn *db.DB, name string, depositAmount int) error {
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

	return nil
}

func seed10Users(ctx context.Context, dbConn *db.DB) error {
	for i := 0; i < 10; i++ {
		SeedUser(ctx, dbConn, "Alex", i*100)
	}

	return nil
}
