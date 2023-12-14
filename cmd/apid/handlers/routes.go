// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package handlers

import (
	"net/http"

	"github.com/timurguseynov/go-wallet-api/internal/platform/db"

	"github.com/timurguseynov/go-wallet-api/internal/middleware"
	"github.com/timurguseynov/go-wallet-api/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(db *db.DB) http.Handler {
	// Create the web handler for setting routes and middleware.
	app := web.New(middleware.RequestLogger, middleware.ErrorHandler)

	// Initialize the routes for the API binding the route to the
	// handler code for each specified verb.

	// user
	u := User{
		MasterDB: db,
	}
	app.Handle(http.MethodPost, "/api/user/create", u.postUserCreate)
	app.Handle(http.MethodPost, "/api/wallet/deposit", u.postUserDeposit)
	app.Handle(http.MethodPost, "/api/wallet/withdraw", u.postUserWithdraw)
	app.Handle(http.MethodGet, "/api/wallet/balance/{userID}", u.getUserBalance)

	// notifier
	n := Notifier{
		MasterDB: db,
	}
	app.Handle("*", "/ws/topic/leaderboard", n.leaderBoard)
	app.Handle("*", "/ws/topic/outcomes", n.outcomes)

	return app
}
