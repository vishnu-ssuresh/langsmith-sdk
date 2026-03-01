package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Request is a transport-level request envelope.
type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Headers http.Header
	Body    []byte
}

// Response is a transport-level response envelope.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// NewRequest creates a request envelope with method and path set.
func NewRequest(method, path string) Request {
	return Request{
		Method: method,
		Path:   path,
	}
}

// WithQuery adds one or more query values to the request.
func (r Request) WithQuery(key string, values ...string) Request {
	if key == "" || len(values) == 0 {
		return r
	}
	if r.Query == nil {
		r.Query = make(url.Values)
	}
	for _, value := range values {
		r.Query.Add(key, value)
	}
	return r
}

// WithHeader sets a request header value.
func (r Request) WithHeader(key, value string) Request {
	if key == "" {
		return r
	}
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	r.Headers.Set(key, value)
	return r
}

// WithBody sets the request body bytes.
func (r Request) WithBody(body []byte) Request {
	r.Body = body
	return r
}

// EncodeJSONBody marshals a value to JSON bytes for Request.Body.
func EncodeJSONBody(value interface{}) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("transport: marshal request body: %w", err)
	}
	return encoded, nil
}
