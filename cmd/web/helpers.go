package main

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug" //able to see stacktrace when there are errors
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) //2 means that if there's an error we want the linenumber and file to be the caller not the callee

	//deal with the error status
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	//can redirect to the error page
	app.clientError(w, http.StatusNotFound)
}

func containsOnlyValidChars(input string) bool {
	validChars := regexp.MustCompile(`^[a-zA-Z ']+$`)
	return validChars.MatchString(input)
}

func (app *application) errRecordNotFound(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	app.clientError(w, http.StatusInternalServerError)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedID")
}
