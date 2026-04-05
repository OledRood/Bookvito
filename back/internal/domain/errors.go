package domain

import "errors"

type ErrorCode string

const (
	ErrorCodeValidation ErrorCode = "validation"
	ErrorCodeForbidden  ErrorCode = "forbidden"
	ErrorCodeNotFound   ErrorCode = "not_found"
	ErrorCodeConflict   ErrorCode = "conflict"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Code)
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewValidationError(message string) error {
	return &AppError{Code: ErrorCodeValidation, Message: message}
}

func NewForbiddenError(message string) error {
	return &AppError{Code: ErrorCodeForbidden, Message: message}
}

func NewNotFoundError(message string) error {
	return &AppError{Code: ErrorCodeNotFound, Message: message}
}

func NewConflictError(message string) error {
	return &AppError{Code: ErrorCodeConflict, Message: message}
}

func AppErrorCode(err error) (ErrorCode, bool) {
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr == nil {
		return "", false
	}
	return appErr.Code, true
}
