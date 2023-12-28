package handlers

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/timurguseynov/go-wallet-api/internal/rest"
	"github.com/timurguseynov/go-wallet-api/internal/user"

	"github.com/pkg/errors"

	"github.com/timurguseynov/go-wallet-api/internal/db"
)

// User represents the User API method handler set.
type User struct {
	MasterDB *db.DB
}

type PostUserAmount struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}

func (a PostUserAmount) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Amount, validation.Required),
		validation.Field(&a.Amount, validation.Min(10)),
	)
}

func (u *User) postUserCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	var userCreate user.User
	err := rest.Unmarshal(r.Body, &userCreate)
	if err != nil {
		return errors.Wrap(err, "")
	}

	id, err := user.Insert(ctx, u.MasterDB, userCreate)
	if err != nil {
		return errors.Wrap(err, "")
	}

	resp := user.User{
		ID: id,
	}

	rest.Respond(ctx, w, resp, http.StatusOK)
	return nil
}

func (u *User) postUserDeposit(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	var userAmount PostUserAmount
	err := rest.Unmarshal(r.Body, &userAmount)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = user.DepositByID(ctx, u.MasterDB, userAmount.ID, userAmount.Amount)
	if err != nil {
		return errors.Wrap(err, "")
	}

	rest.Respond(ctx, w, true, http.StatusOK)
	return nil
}

func (u *User) postUserWithdraw(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	var userAmount PostUserAmount
	err := rest.Unmarshal(r.Body, &userAmount)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = user.WithdrawByID(ctx, u.MasterDB, userAmount.ID, userAmount.Amount)
	if err != nil {
		if err == user.ErrInsufficientFunds {
			return rest.NewResponseError(err, http.StatusPaymentRequired)
		}
		return errors.Wrap(err, "")
	}

	rest.Respond(ctx, w, true, http.StatusOK)
	return nil
}

func (u *User) getUserBalance(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	userBalance, err := user.GetBalanceByID(ctx, u.MasterDB, params["userID"])
	if err != nil {
		return errors.Wrap(err, "")
	}

	b := user.User{
		Balance: userBalance,
	}

	rest.Respond(ctx, w, b, http.StatusOK)
	return nil
}
