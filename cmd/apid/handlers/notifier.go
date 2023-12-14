package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/timurguseynov/user-wallet-api/internal/platform/db"
	"github.com/timurguseynov/user-wallet-api/internal/platform/user"
)

type Notifier struct {
	MasterDB *db.DB
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (n *Notifier) leaderBoard(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	wsConn, _ := upgrader.Upgrade(w, r, nil)

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			users, err := user.List(ctx, n.MasterDB)
			if err != nil {
				return errors.Wrap(err, "")
			}

			// put leaders first
			sort.Slice(users, func(i, j int) bool {
				return users[i].Balance > users[j].Balance
			})

			message, err := json.Marshal(users)
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = wsConn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}
	}
}

func (n *Notifier) outcomes(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	wsConn, _ := upgrader.Upgrade(w, r, nil)

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			users, err := user.List(ctx, n.MasterDB)
			if err != nil {
				return errors.Wrap(err, "")
			}

			message, err := json.Marshal(users)
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = wsConn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}
	}
}
