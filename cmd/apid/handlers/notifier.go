package handlers

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/timurguseynov/go-wallet-api/internal/platform/db"
	"github.com/timurguseynov/go-wallet-api/internal/platform/rest"
	"github.com/timurguseynov/go-wallet-api/internal/platform/user"
)

type Notifier struct {
	MasterDB *db.DB
}

func (n *Notifier) leaderBoard(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastCheckUsers []user.User

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			users, err := user.ListLeaders(ctx, n.MasterDB)
			if err != nil {
				return errors.Wrap(err, "")
			}

			// only send data if it's changed
			if reflect.DeepEqual(lastCheckUsers, users) {
				continue
			}
			lastCheckUsers = users

			err = rest.WebsocketRespond(ctx, users)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}
	}
}

func (n *Notifier) outcomes(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastCheckUsers []user.User

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			users, err := user.List(ctx, n.MasterDB)
			if err != nil {
				return errors.Wrap(err, "")
			}

			// only send data if it's changed
			if reflect.DeepEqual(lastCheckUsers, users) {
				continue
			}
			lastCheckUsers = users

			err = rest.WebsocketRespond(ctx, users)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}
	}
}
