package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/timurguseynov/go-wallet-api/cmd/apid/handlers"
	"github.com/timurguseynov/go-wallet-api/internal/platform/user"
	"github.com/timurguseynov/go-wallet-api/internal/platform/web"
)

var (
	userID         string
	depositAmount  int64 = 10000
	withdrawAmount int64 = 5000
)

func RunTestUser(t *testing.T) {
	t.Run("postUserCreate", postUserCreate)
	t.Run("postUserDeposit", postUserDeposit)
	t.Run("postUserDepositValidateAmount", postUserDepositValidateAmount)
	t.Run("postUserWithdraw", postUserWithdraw)
	t.Run("postUserWithdrawValidateAmount", postUserWithdrawValidateAmount)
	t.Run("getUserBalance", getUserBalance)
}

func postUserCreate(t *testing.T) {
	expected := user.User{
		Name: "Alex",
	}
	body, err := json.Marshal(expected)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))
	var got user.User
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	require.NotEqual(t, "", got)

	userID = got.ID
}

func postUserDeposit(t *testing.T) {
	expected := handlers.PostUserAmount{
		ID:     userID,
		Amount: depositAmount,
	}
	body, err := json.Marshal(expected)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/deposit", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))
	var got bool
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	require.True(t, got)
}

func postUserDepositValidateAmount(t *testing.T) {
	expected := handlers.PostUserAmount{
		ID:     userID,
		Amount: 0,
	}
	body, err := json.Marshal(expected)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/deposit", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code, http.StatusText(w.Code))

	var got web.JSONError
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, web.ErrValidation.Error(), got.Error)
	require.Equal(t, "amount", got.Fields[0].Fld)
	require.Equal(t, "cannot be blank", got.Fields[0].Err)
}

func postUserWithdraw(t *testing.T) {
	expected := handlers.PostUserAmount{
		ID:     userID,
		Amount: withdrawAmount,
	}
	body, err := json.Marshal(expected)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/withdraw", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))
	var got bool
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	require.True(t, got)
}

func postUserWithdrawValidateAmount(t *testing.T) {
	expected := handlers.PostUserAmount{
		ID:     userID,
		Amount: 0,
	}
	body, err := json.Marshal(expected)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/withdraw", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code, http.StatusText(w.Code))

	var got web.JSONError
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, web.ErrValidation.Error(), got.Error)
	require.Equal(t, "amount", got.Fields[0].Fld)
	require.Equal(t, "cannot be blank", got.Fields[0].Err)
}

func getUserBalance(t *testing.T) {
	t.Log(fmt.Sprintf("/api/wallet/balance/%s", userID))
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/wallet/balance/%s", userID), nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))

	var got user.User
	err := json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, depositAmount-withdrawAmount, got.Balance)
}
