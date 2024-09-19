package errs

import (
	"fmt"

	"github.com/FlyKarlik/auth-service/pkg/codes"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func GetCodeFromError(err error) int {
	customErr, _ := err.(*Error)
	return customErr.Code
}

func GetMessageFromError(err error) string {
	customErr, _ := err.(*Error)
	return customErr.Message
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func Newf(code int, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

var (
	ErrServiceNameNotConfigured   = New(codes.ErrorInternal, "service name field is not configured")
	ErrServerHostNotConfigured    = New(codes.ErrorInternal, "server host name field is not configured")
	ErrDatabaseUrlNotConfigured   = New(codes.ErrorInternal, "database url field is not configured")
	ErrJaegerHostNotConfigured    = New(codes.ErrorInternal, "jaeger host field is not configured")
	ErrLogLevelNotConfigured      = New(codes.ErrorInternal, "log level field is not configured")
	ErrJwtSecreNotConfigured      = New(codes.ErrorInternal, "jwt secret field is not configured")
	ErrClientIPNotFound           = New(codes.ErrorInternal, "client ip not found")
	ErrInvalidClientIP            = New(codes.ErrorInternal, "invalid client ip")
	ErrCreateJWT                  = New(codes.ErrorInternal, "cannot create jwt token")
	ErrInvalidToken               = New(codes.ErrorUnauthorized, "invalid token")
	ErrMismatchUserData           = New(codes.ErrorUnauthorized, "mismatch user id or client ip")
	ErrMismatchTokenVariety       = New(codes.ErrorUnauthorized, "mismatch token variety")
	ErrDatabaseExecContext        = New(codes.ErrorInternal, "database exec context error")
	ErrDatabaseGetContext         = New(codes.ErrorInternal, "database get context error")
	ErrUserIDNotFound             = New(codes.ErrorBadRequest, "user id param not found")
	ErrInvalidUserID              = New(codes.ErrorBadRequest, "invalid user id")
	ErrIncorrectRerfreshOperation = New(codes.ErrorUnauthorized, "incorrect refresh operation")
)
