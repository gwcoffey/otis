package oerr

import (
	"errors"
)

type ErrorCode int

type OtisError struct {
	Code    ErrorCode
	Message string
}

const (
	projectNotFound ErrorCode = iota + 1
)

func ProjectNotFound() *OtisError {
	return &OtisError{Code: projectNotFound, Message: "this is not an otis directory"}
}

func IsProjectNotFoundErr(err error) bool {
	var otisError *OtisError
	if errors.As(err, &otisError) {
		return otisError.Code == projectNotFound
	}
	return false
}

func (e *OtisError) Error() string {
	return e.Message
}
