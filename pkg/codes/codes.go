package codes

import "net/http"

const (
	StatusOK          = http.StatusOK
	ErrorUnauthorized = http.StatusUnauthorized
	ErrorBadRequest   = http.StatusBadRequest
	ErrorNotFound     = http.StatusNotFound
	ErrorInternal     = http.StatusInternalServerError
)

const (
	Success = "success"
	Failure = "failure"
)
