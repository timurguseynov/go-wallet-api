package tests

import (
	"context"
	"encoding/json"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"github.com/timurguseynov/go-wallet-api/internal/platform/user"
)

func RunTestNotifier(t *testing.T) {
	err := addTestUsers(context.Background(), test.MasterDB)
	require.NoError(t, err)

	t.Run("wsNotifierLeaderBoard", wsNotifierLeaderBoard)
	// t.Run("wsNotifierOutcomes", wsNotifierOutcomes)
}

func wsNotifierLeaderBoard(t *testing.T) {
	s := httptest.NewServer(a)
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1) + "/ws/topic/leaderboard"
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.NoError(t, err)
	defer ws.Close()

	messageType, message, err := ws.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, messageType)

	var users []user.User
	err = json.Unmarshal(message, &users)
	require.NoError(t, err)
	require.Equal(t, int64(5000), users[0].Balance)
	require.Equal(t, int64(900), users[1].Balance)
	require.Equal(t, int64(800), users[2].Balance)
}

func wsNotifierOutcomes(t *testing.T) {
	s := httptest.NewServer(a)
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1) + "/ws/topic/outcomes"
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.NoError(t, err)
	defer ws.Close()

	messageType, message, err := ws.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, messageType)

	var users []user.User
	err = json.Unmarshal(message, &users)
	require.NoError(t, err)
	require.Equal(t, 11, len(users))
}
