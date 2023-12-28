package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timurguseynov/go-wallet-api/cmd/apid/handlers"
	"github.com/timurguseynov/go-wallet-api/internal/rest"
	"github.com/timurguseynov/go-wallet-api/internal/user"
)

var (
	userID         string
	depositAmount  int64 = 10000
	withdrawAmount int64 = 5000
)

func RunTestUser(t *testing.T) {
	t.Run("postUserCreate", postUserCreate)
	t.Run("postUserDeposit", postUserDeposit)
	t.Run("postUserDepositValidateAmount", postUserDepositValidateInputAmount)
	t.Run("postUserWithdraw", postUserWithdraw)
	t.Run("postUserWithdrawInsufficientFunds", postUserWithdrawInsufficientFunds)
	t.Run("postUserWithdrawValidateAmount", postUserWithdrawValidateInputAmount)
	t.Run("getUserBalance", getUserBalance)
}

func postUserCreate(t *testing.T) {
	u := user.User{
		Name: "Alex",
	}
	body, err := json.Marshal(u)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))
	var got user.User
	err = json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.NotEqual(t, "", got)
	userID = got.ID
}

func postUserDeposit(t *testing.T) {
	userAmount := handlers.PostUserAmount{
		ID:     userID,
		Amount: depositAmount,
	}
	body, err := json.Marshal(userAmount)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/deposit", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))
	var got bool
	err = json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.True(t, got)
}

func postUserDepositValidateInputAmount(t *testing.T) {
	userAmount := handlers.PostUserAmount{
		ID:     userID,
		Amount: 0,
	}
	body, err := json.Marshal(userAmount)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/deposit", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code, http.StatusText(w.Code))

	var got rest.JSONError
	err = json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, rest.ErrValidation.Error(), got.Error)
	assert.Equal(t, "amount", got.Fields[0].Fld)
	assert.Equal(t, "cannot be blank", got.Fields[0].Err)
}

func postUserWithdraw(t *testing.T) {
	userAmount := handlers.PostUserAmount{
		ID:     userID,
		Amount: withdrawAmount,
	}
	body, err := json.Marshal(userAmount)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/withdraw", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))
	var got bool
	err = json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.True(t, got)
}

func postUserWithdrawInsufficientFunds(t *testing.T) {
	userAmount := handlers.PostUserAmount{
		ID:     userID,
		Amount: withdrawAmount + 1,
	}
	body, err := json.Marshal(userAmount)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/withdraw", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusPaymentRequired, w.Code, http.StatusText(w.Code))
	var got rest.JSONError
	err = json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, user.ErrInsufficientFunds.Error(), got.Error)
}

func postUserWithdrawValidateInputAmount(t *testing.T) {
	userAmount := handlers.PostUserAmount{
		ID:     userID,
		Amount: 0,
	}
	body, err := json.Marshal(userAmount)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/wallet/withdraw", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code, http.StatusText(w.Code))

	var got rest.JSONError
	err = json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, rest.ErrValidation.Error(), got.Error)
	assert.Equal(t, "amount", got.Fields[0].Fld)
	assert.Equal(t, "cannot be blank", got.Fields[0].Err)
}

func getUserBalance(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/wallet/balance/%s", userID), nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, http.StatusText(w.Code))

	var got user.User
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, depositAmount-withdrawAmount, got.Balance)
}
