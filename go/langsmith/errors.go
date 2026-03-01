package langsmith

import (
	"errors"
	"net/http"
)

var (
	// ErrInvalidConfig means required SDK config is missing or malformed.
	ErrInvalidConfig = errors.New("langsmith: invalid config")
	// ErrUnauthorized means credentials are invalid or missing permissions.
	ErrUnauthorized = errors.New("langsmith: unauthorized")
	// ErrForbidden means the caller is authenticated but not allowed.
	ErrForbidden = errors.New("langsmith: forbidden")
	// ErrNotFound means the requested resource does not exist.
	ErrNotFound = errors.New("langsmith: not found")
	// ErrRateLimited means requests are being throttled.
	ErrRateLimited = errors.New("langsmith: rate limited")
	// ErrTransient means the request failed due to a retryable/server-side issue.
	ErrTransient = errors.New("langsmith: transient error")
)

// ErrorForStatus maps HTTP response status codes to SDK sentinel errors.
func ErrorForStatus(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return ErrTransient
	default:
		return nil
	}
}
