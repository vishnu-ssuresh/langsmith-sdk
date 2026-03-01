package langsmith

import (
	"net/http"
	"testing"
)

func TestErrorForStatus(t *testing.T) {
	tests := []struct {
		status int
		want   error
	}{
		{status: http.StatusUnauthorized, want: ErrUnauthorized},
		{status: http.StatusForbidden, want: ErrForbidden},
		{status: http.StatusNotFound, want: ErrNotFound},
		{status: http.StatusTooManyRequests, want: ErrRateLimited},
		{status: http.StatusInternalServerError, want: ErrTransient},
		{status: http.StatusBadGateway, want: ErrTransient},
		{status: http.StatusServiceUnavailable, want: ErrTransient},
		{status: http.StatusGatewayTimeout, want: ErrTransient},
		{status: http.StatusBadRequest, want: nil},
	}

	for _, tt := range tests {
		if got := ErrorForStatus(tt.status); got != tt.want {
			t.Fatalf("ErrorForStatus(%d) = %v, want %v", tt.status, got, tt.want)
		}
	}
}
