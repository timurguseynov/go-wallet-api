// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Current Status Codes:
//		200 OK           : StatusOK                  : Call is success and returning data.
//		204 No Content   : StatusNoContent           : Call is success and returns no data.
//		400 Bad Request  : StatusBadRequest          : Invalid post data (syntax or semantics).
//		401 Unauthorized : StatusUnauthorized        : Authentication failure.
//		404 Not Found    : StatusNotFound            : Invalid URL or identifier.
//		500 Internal     : StatusInternalServerError : Application specific beyond scope of user.

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var ErrNoWebsocketConnection = errors.New("no websocket connection found")

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string      `json:"field_name"`
	Err interface{} `json:"error"`
}

// InvalidError is a custom error type for invalid fields.
type InvalidError []Invalid

// Error implements the error interface for InvalidError.
func (err InvalidError) Error() string {
	str := "{"
	for _, v := range err {
		str += fmt.Sprintf("%s: %s,", v.Fld, v.Err)
	}
	str = strings.TrimRight(str, ",")
	str += "}"
	return str
}

// ResponseError is used to pass an error during the request through the
// application with web specific context.
type ResponseError struct {
	Err    error
	Status int
}

// NewResponseError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewResponseError(err error, status int) error {
	return ResponseError{err, status}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (re ResponseError) Error() string {
	return re.Err.Error()
}

// JSONError is the response for errors that occur within the API.
type JSONError struct {
	Error  string       `json:"error"`
	Fields InvalidError `json:"fields,omitempty"`
}

var (
	// ErrUnauthorized occurs when the call is not authorized.
	ErrUnauthorized = errors.New("Not authorized")

	// ErrDBNotConfigured occurs when the DB is not initialized.
	ErrDBNotConfigured = errors.New("DB not initialized")

	// ErrNotFound is abstracting the mgo not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in it's proper form")

	// ErrValidation occurs when there are validation errors.
	ErrValidation = errors.New("Validation errors occurred")

	// ErrForbidden occurs when we know who the user is but they attempt a
	// forbidden action.
	ErrForbidden = errors.New("Forbidden")
)

// ErrorHandler handles all error responses for the API.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, err error) {
	switch errors.Cause(err) {
	case ErrNotFound:
		RespondError(ctx, w, err, http.StatusNotFound)
		return

	case ErrInvalidID:
		RespondError(ctx, w, err, http.StatusBadRequest)
		return

	case ErrValidation:
		RespondError(ctx, w, err, http.StatusBadRequest)
		return

	case ErrUnauthorized:
		RespondError(ctx, w, err, http.StatusUnauthorized)
		return

	case ErrForbidden:
		RespondError(ctx, w, err, http.StatusForbidden)
		return
	}

	switch e := errors.Cause(err).(type) {
	case InvalidError:
		v := JSONError{
			Error:  ErrValidation.Error(),
			Fields: e,
		}
		Respond(ctx, w, v, http.StatusBadRequest)
		return
	case ResponseError:
		RespondError(ctx, w, e.Err, e.Status)
		return
	}

	RespondError(ctx, w, err, http.StatusInternalServerError)
}

// RespondError sends JSON describing the error
func RespondError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	Respond(ctx, w, JSONError{Error: err.Error()}, code)
}

// Respond sends JSON to the client.
// If code is StatusNoContent, v is expected to be nil.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, code int) {
	// Set the status code for the request logger middleware.
	v := ctx.Value(KeyValues).(*Values)
	v.StatusCode = code

	// Just set the status code and we are done.
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	// Set the content type.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response and context.
	w.WriteHeader(code)
	// Marshal the data into a JSON string.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logStdErr.Printf("%s : Respond %v Marshalling JSON response\n", v.TraceID, err)
		jsonData = []byte("{}")
	}

	// Send the result back to the client.
	io.WriteString(w, string(jsonData))
}

func WebsocketErrorHandler(ctx context.Context, err error) {
	// TODO: replace status codes

	serverErrorStatusCode := 1011

	switch errors.Cause(err) {
	case ErrNotFound:
		WebsocketRespondError(ctx, err, serverErrorStatusCode)
		return

	case ErrInvalidID:
		WebsocketRespondError(ctx, err, serverErrorStatusCode)
		return

	case ErrValidation:
		WebsocketRespondError(ctx, err, serverErrorStatusCode)
		return

	case ErrUnauthorized:
		WebsocketRespondError(ctx, err, serverErrorStatusCode)
		return

	case ErrForbidden:
		WebsocketRespondError(ctx, err, serverErrorStatusCode)
		return
	}

	switch e := errors.Cause(err).(type) {
	case InvalidError:
		v := JSONError{
			Error:  ErrValidation.Error(),
			Fields: e,
		}
		WebsocketRespond(ctx, v)
		return
	case ResponseError:
		WebsocketRespondError(ctx, e.Err, serverErrorStatusCode)
		return
	}

	WebsocketRespondError(ctx, err, serverErrorStatusCode)
}

func isWebsocket(ctx context.Context) bool {
	_, ok := ctx.Value(WebsocketConnection).(*websocket.Conn)
	return ok
}

func WebsocketRespond(ctx context.Context, data interface{}) error {
	wsConn, ok := ctx.Value(WebsocketConnection).(*websocket.Conn)
	if !ok {
		return ErrNoWebsocketConnection
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = wsConn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func WebsocketRespondError(ctx context.Context, data interface{}, code int) {
	wsConn, ok := ctx.Value(WebsocketConnection).(*websocket.Conn)
	if !ok {
		logStdErr.Println(ErrNoWebsocketConnection)
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logStdErr.Println(err)
		return
	}

	message := websocket.FormatCloseMessage(code, string(jsonData))

	err = wsConn.WriteMessage(websocket.CloseMessage, message)
	if err != nil {
		logStdErr.Println(err)
		return
	}

	wsConn.Close()

}
