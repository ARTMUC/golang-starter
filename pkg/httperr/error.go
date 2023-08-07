package httperr

import "net/http"

type ErrCustomError struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode"`
	Err        error  `json:"-"`
}

func (e ErrCustomError) Error() string {
	return e.Message
}

func NewHttpError(msg string, statusCode int, err error) ErrCustomError {
	return ErrCustomError{
		Message:    msg,
		StatusCode: statusCode,
		Err:        err,
	}
}

func NewNotFoundError(msg string, err error) ErrCustomError {
	return NewHttpError(msg, http.StatusNotFound, err)
}

func NewBadRequestError(msg string, err error) ErrCustomError {
	return NewHttpError(msg, http.StatusBadRequest, err)
}

func NewUnAuthorizedError(msg string, err error) ErrCustomError {
	return NewHttpError(msg, http.StatusUnauthorized, err)
}
