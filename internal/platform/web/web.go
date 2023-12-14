// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Package web provides a thin layer of support for writing web services.
package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"

	"github.com/pborman/uuid"
)

// TraceIDHeader is the header added to outgoing requests which adds the
// traceID to it.
const TraceIDHeader = "X-Trace-ID"

// Unmarshal decodes the input to the struct type and checks the
// fields to verify the value is in a proper state.
func Unmarshal(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	return validate(v)
}

func validate(v interface{}) error {
	// check if it's validatable
	validatable, ok := v.(validation.Validatable)
	if !ok {
		return nil
	}

	// validate
	err := validatable.Validate()

	// check for errors
	validationErr, ok := err.(validation.Errors)
	if !ok {
		return nil
	}

	// format errors
	var inv InvalidError
	for key, err := range validationErr {
		var e interface{}
		if ms, ok := err.(json.Marshaler); ok {
			e = ms
		} else {
			e = err.Error()
		}

		inv = append(inv, Invalid{Fld: key, Err: e})
	}

	return inv
}

// Key represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

// A Middleware is a type that wraps a handler to remove boilerplate or other
// concerns not direct to any given Handler.
type Middleware func(Handler) Handler

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*mux.Router
	mw []Middleware
}

// New creates an App value that handle a set of routes for the application.
// You can provide any number of middleware and they'll be used to wrap every
// request handler.
func New(mw ...Middleware) *App {
	return &App{
		Router: mux.NewRouter(),
		mw:     mw,
	}
}

// Use adds the set of provided middleware onto the Application middleware
// chain. Any route running off of this App will use all the middleware provided
// this way always regardless of the ordering of the Handle/Use functions.
func (a *App) Use(mw ...Middleware) {
	a.mw = append(a.mw, mw...)
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// Wrap up the application-wide first, this will call the first function
	// of each middleware which will return a function of type Handler. Each
	// Handler will then be wrapped up with the other handlers from the chain.
	handler = wrapMiddleware(wrapMiddleware(handler, mw), a.mw)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: uuid.New(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		// Set the trace id on the outgoing requests before any other header to
		// ensure that the trace id is ALWAYS added to the request regardless of
		// any error occuring or not.
		w.Header().Set(TraceIDHeader, v.TraceID)

		// Extract url params like /product/:id to use as a map in handler
		vars := mux.Vars(r)

		// Call the wrapped handler functions.
		handler(ctx, w, r, vars)
	}

	// Add this handler for the specified verb and route.
	a.Router.HandleFunc(path, h).Methods(verb)
}

func (a *App) AddHandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	a.Router.HandleFunc(path, f)
}

// wrapMiddleware wraps a handler with some middleware.
func wrapMiddleware(handler Handler, mw []Middleware) Handler {

	// Wrap with our group specific middleware.
	for i := len(mw) - 1; i >= 0; i-- {
		if mw[i] != nil {
			handler = mw[i](handler)
		}
	}

	return handler
}
