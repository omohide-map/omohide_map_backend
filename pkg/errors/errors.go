package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) HTTPStatus() int {
	return e.Code
}

var (
	ErrBadRequest          = &AppError{Code: http.StatusBadRequest, Message: "Bad request"}
	ErrUnauthorized        = &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden           = &AppError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrNotFound            = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrConflict            = &AppError{Code: http.StatusConflict, Message: "Resource conflict"}
	ErrInternalServer      = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrServiceUnavailable  = &AppError{Code: http.StatusServiceUnavailable, Message: "Service unavailable"}
	ErrUnprocessableEntity = &AppError{Code: http.StatusUnprocessableEntity, Message: "Unprocessable entity"}
)

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewWithDetail(code int, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func Wrap(err error, appErr *AppError) *AppError {
	return &AppError{
		Code:    appErr.Code,
		Message: appErr.Message,
		Detail:  appErr.Detail,
		Err:     err,
	}
}

func WrapWithMessage(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func InvalidRequest(detail string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid request",
		Detail:  detail,
	}
}

func ValidationError(detail string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Validation error",
		Detail:  detail,
	}
}

func AuthenticationRequired() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Authentication required",
	}
}

func InvalidToken(detail string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid token",
		Detail:  detail,
	}
}

func MissingAuthHeader() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Missing authorization header",
	}
}

func InvalidAuthFormat() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid authorization format",
	}
}

func UserIDNotFound() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "User ID not found",
	}
}

func ResourceNotFound(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func DatabaseError(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Database operation failed",
		Err:     err,
	}
}

func StorageError(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Storage operation failed",
		Err:     err,
	}
}

func EnvironmentVariableError(varName string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Configuration error",
		Detail:  fmt.Sprintf("%s environment variable is required", varName),
	}
}

func ImageProcessingError(err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Image processing failed",
		Err:     err,
	}
}

func RequestTimeNotFound() *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Request time not found",
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return WrapWithMessage(err, http.StatusInternalServerError, "An unexpected error occurred")
}
