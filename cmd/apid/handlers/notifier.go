package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/timurguseynov/go-wallet-api/internal/platform/db"
	"github.com/timurguseynov/go-wallet-api/internal/platform/user"
)

type Notifier struct {
	MasterDB *db.DB
}

type NotifierError struct {
	Error string `json:"error"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (n *Notifier) leaderBoard(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	wsConn, _ := upgrader.Upgrade(w, r, nil)

	ticker := time.NewTicker(1 * time.Second)

	var lastCheckUsers []user.User

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			users, err := user.List(ctx, n.MasterDB)
			if err != nil {
				sendMessage(wsConn, NotifierError{"couldn't find users"})
				wsConn.Close()
				return errors.Wrap(err, "")
			}

			// put leaders first
			sort.Slice(users, func(i, j int) bool {
				return users[i].Balance > users[j].Balance
			})

			// only send data if it's changed
			if reflect.DeepEqual(lastCheckUsers, users) {
				continue
			}
			lastCheckUsers = users

			err = sendMessage(wsConn, users)
			if err != nil {
				wsConn.Close()
				return errors.Wrap(err, "")
			}
		}
	}
}

func (n *Notifier) outcomes(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	wsConn, _ := upgrader.Upgrade(w, r, nil)

	ticker := time.NewTicker(1 * time.Second)

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

			message, err := json.Marshal(users)
			if err != nil {
				wsConn.Close()
				return errors.Wrap(err, "")
			}

			err = wsConn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				wsConn.Close()
				return errors.Wrap(err, "")
			}
		}
	}
}

func sendMessage(wsConn *websocket.Conn, messageStruct interface{}) error {
	message, err := json.Marshal(messageStruct)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = wsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
