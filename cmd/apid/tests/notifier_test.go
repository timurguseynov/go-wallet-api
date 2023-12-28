package tests

import (
	"context"
	"encoding/json"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/timurguseynov/go-wallet-api/internal/tests"
	"github.com/timurguseynov/go-wallet-api/internal/user"
)

func RunTestNotifier(t *testing.T) {
	t.Run("wsNotifierLeaderBoard", wsNotifierLeaderBoard)
	t.Run("wsNotifierOutcomes", wsNotifierOutcomes)
}

func wsNotifierLeaderBoard(t *testing.T) {
	s := httptest.NewServer(a)
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1) + "/ws/topic/leaderboard"
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer ws.Close()

	messageType, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, websocket.TextMessage, messageType)

	var users []user.User
	err = json.Unmarshal(message, &users)
	assert.NoError(t, err)
	assert.True(t, len(users) > 2)

	// change data to allow one more read
	err = tests.SeedUser(context.TODO(), test.MasterDB, "John1", 100)
	assert.NoError(t, err)

	// test one more read
	messageType, message, err = ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, websocket.TextMessage, messageType)
	err = json.Unmarshal(message, &users)
	assert.NoError(t, err)
	assert.True(t, len(users) > 2)

}

func wsNotifierOutcomes(t *testing.T) {
	s := httptest.NewServer(a)
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1) + "/ws/topic/outcomes"
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer ws.Close()

	messageType, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, websocket.TextMessage, messageType)

	var users []user.User
	err = json.Unmarshal(message, &users)
	assert.NoError(t, err)
	assert.True(t, len(users) > 2)

	// change data to allow one more read
	err = tests.SeedUser(context.TODO(), test.MasterDB, "John1", 100)
	assert.NoError(t, err)

	// test one more read
	messageType, message, err = ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, websocket.TextMessage, messageType)
}
