// Package main provides various handlers for handling user requests.
//
// The functions in this file handle user requests for login, registration, video upload, video termination,
// video display, video editing, video deletion, and other related operations.
// The handlers interact with the database and the user interface to provide the required functionality.
//
// Global Variables:
// Data is a global variable of type map with string keys and values of any type (interface{}).
// It is used to hold various data, such as video information, that needs to be accessed across different functions.
package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError is a function that handles server errors.
// It takes a http.ResponseWriter and an error as arguments.
// The function first creates a stack trace for the error, which includes the error message and the stack trace from the point where the error occurred.
// This stack trace is then logged using the application's error logger.
// Finally, the function sends a HTTP 500 Internal Server Error response to the client.
// This function is typically used in other handlers to handle unexpected errors and provide a consistent error response.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
