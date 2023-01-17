package errors

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	Unknown ErrorCode = iota
	NotFound
	JsonError
)

// The Service error can be used when the type of error from a lower level needs to be checked.
// e.g. a service func may return an service error with an ErrorCode of notfound or unknown
// the caller can then decide at the API boundary what EEMI errors need to be returned.
// This approach lets the decision on the eemi error code be made at the the last responsible moment
type Error struct {
	error
	ErrorCode
}

func (e Error) Error() string {
	return fmt.Sprintf("%+v", e.error)
}

func NewError(err error, errCode ErrorCode) error {
	return Error{error: err, ErrorCode: errCode}
}

func NewErrorStr(err string, errCode ErrorCode) error {
	return Error{error: errors.New(err), ErrorCode: errCode}
}
