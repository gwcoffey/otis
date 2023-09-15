package oerr

import (
	"errors"
	"fmt"
)

type ErrorCode int

type OtisError struct {
	Code    ErrorCode
	Message string
}

const (
	projectNotFound ErrorCode = iota + 1
	folderPathNotFound
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

func FolderPathNotFound(path string) *OtisError {
	return &OtisError{Code: folderPathNotFound, Message: fmt.Sprintf("path %s does not exist or is not a folder", path)}
}

func (e *OtisError) Error() string {
	return e.Message
}
