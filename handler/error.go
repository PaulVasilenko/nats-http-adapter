package handler

import (
	"net/http"

	"github.com/pkg/errors"
)

type HTTPError struct {
	parent  error
	message string
	code    int
}

func (e HTTPError) Error() string {
	if e.parent != nil {
		return errors.Wrap(e.parent, e.message).Error()
	}
	return e.message
}

func (e HTTPError) Code() int {
	return e.code
}

func BadRequest(message string) HTTPError {
	return HTTPError{
		message: message,
		code:    http.StatusBadRequest,
	}
}
