package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
)

var logStdErr *log.Logger

func init() {
	logStdErr = log.New(os.Stderr, "", log.LstdFlags)
	logStdErr.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// ErrorHandler for catching and responding errors.
func ErrorHandlerMiddleware(next Handler) Handler {
	// Create the handler that will be attached in the middleware chain.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

		v := ctx.Value(KeyValues).(*Values)

		// In the event of a panic, we want to capture it here so we can send an
		// error down the stack.
		defer func() {

			if r := recover(); r != nil {

				// Log the panic.
				logStdErr.Printf("%s : ERROR : Panic Caught : %s\n", v.TraceID, r)

				// Respond with the error.
				RespondError(ctx, w, errors.New("unhandled"), http.StatusInternalServerError)

				// Print out the stack.
				logStdErr.Printf("%s : ERROR : Stacktrace\n%s\n", v.TraceID, debug.Stack())
			}
		}()

		// TODO: check that no patient sensitive information is leaked
		if err := next(ctx, w, r, params); err != nil {
			if errors.Cause(err) != ErrNotFound {

				// Log the error.
				logStdErr.Printf("%s : ERROR : %+v\n", v.TraceID, err)
			}

			// Respond with the error.
			if isWebsocket(ctx) {
				WebsocketErrorHandler(ctx, err)
				return nil
			}

			ErrorHandler(ctx, w, err)

			// The error has been handled so we can stop propigating it.
			return nil
		}

		return nil
	}
}

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLoggerMiddleware(next Handler) Handler {
	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(KeyValues).(*Values)

		next(ctx, w, r, params)

		log.Printf("%s : (%d) : %s %s -> %s (%s)",
			v.TraceID,
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Now),
		)

		// This is the top of the food chain. At this point all error
		// handling has been done including logging.
		return nil
	}
}

func websocketMiddleware(next Handler) Handler {
	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		// prevent double upgrade
		if isWebsocket(ctx) {
			return next(ctx, w, r, params)
		}

		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return errors.Wrap(err, "")
		}

		ctx = context.WithValue(ctx, WebsocketConnection, wsConn)

		return next(ctx, w, r, params)
	}
}
