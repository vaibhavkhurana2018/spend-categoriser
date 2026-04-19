package handler

import "net/http"

// httpError is a custom error type that carries an HTTP status code.
// GoFr uses the StatusCode() method to determine the response status.
type httpError struct {
	code    int
	message string
}

func (e httpError) Error() string   { return e.message }
func (e httpError) StatusCode() int { return e.code }

func errBadRequest(msg string) httpError {
	return httpError{code: http.StatusBadRequest, message: msg}
}

func errEntityTooLarge(msg string) httpError {
	return httpError{code: http.StatusRequestEntityTooLarge, message: msg}
}

func errUnprocessable(msg string) httpError {
	return httpError{code: http.StatusUnprocessableEntity, message: msg}
}
