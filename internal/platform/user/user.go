package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/timurguseynov/user-wallet-api/internal/platform/db"
)

type User struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Balance int64  `json:"balance,omitempty"`
}

func Insert(ctx context.Context, dbConn *db.DB, u User) (string, error) {
	txn := dbConn.Txn(true)

	u.ID = uuid.New().String()

	if err := txn.Insert("user", u); err != nil {
		return "", errors.Wrap(err, "committing transaction")
	}

	txn.Commit()

	return u.ID, nil
}

func GetByID(ctx context.Context, dbConn *db.DB, userID string) (*User, error) {
	txn := dbConn.Txn(true)

	raw, err := txn.First("user", "id", userID)
	if err != nil {
		return nil, errors.Wrap(err, "txn.First")
	}

	user, ok := raw.(User)
	if !ok {
		return nil, errors.New("couldn't type assert user")
	}

	txn.Commit()

	return &user, nil
}

func DepositByID(ctx context.Context, dbConn *db.DB, userID string, amount int64) error {
	txn := dbConn.Txn(true)

	raw, err := txn.First("user", "id", userID)
	if err != nil {
		return errors.Wrap(err, "txn.First")
	}

	user, ok := raw.(User)
	if !ok {
		return errors.New("couldn't type assert user")
	}

	user.Balance = user.Balance + amount

	if err := txn.Insert("user", user); err != nil {
		return errors.Wrap(err, "txn.Insert")
	}

	txn.Commit()

	return nil
}

func WithdrawByID(ctx context.Context, dbConn *db.DB, userID string, amount int64) error {

	txn := dbConn.Txn(true)

	raw, err := txn.First("user", "id", userID)
	if err != nil {
		return errors.Wrap(err, "txn.First")
	}

	user, ok := raw.(User)
	if !ok {
		return errors.New("couldn't type assert user")
	}

	user.Balance = user.Balance - amount

	if err := txn.Insert("user", user); err != nil {
		return errors.Wrap(err, "txn.Insert")
	}

	txn.Commit()

	return nil
}

func GetBalanceByID(ctx context.Context, dbConn *db.DB, userID string) (int64, error) {
	txn := dbConn.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("user", "id", userID)
	if err != nil {
		return 0, errors.Wrap(err, "txn.First")
	}

	user, ok := raw.(User)
	if !ok {
		return 0, errors.New("couldn't type assert user")
	}

	return user.Balance, nil
}

func List(ctx context.Context, dbConn *db.DB) ([]User, error) {
	var users []User

	txn := dbConn.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("user", "id")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		u, ok := obj.(User)
		if !ok {
			return nil, errors.New("couldn't type assert user")
		}
		users = append(users, u)
	}

	return users, nil
}
