package transport

import (
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest("POST", "/runs")
	if req.Method != "POST" {
		t.Fatalf("Method = %q, want %q", req.Method, "POST")
	}
	if req.Path != "/runs" {
		t.Fatalf("Path = %q, want %q", req.Path, "/runs")
	}
}

func TestRequestBuilders(t *testing.T) {
	req := NewRequest("GET", "/runs").
		WithQuery("tag", "a", "b").
		WithHeader("X-Test", "yes").
		WithBody([]byte(`{"ok":true}`))

	if got := req.Query["tag"]; len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("Query[tag] = %#v, want [a b]", got)
	}
	if got := req.Headers.Get("X-Test"); got != "yes" {
		t.Fatalf("Header X-Test = %q, want %q", got, "yes")
	}
	if got := string(req.Body); got != `{"ok":true}` {
		t.Fatalf("Body = %q, want %q", got, `{"ok":true}`)
	}
}

func TestEncodeJSONBody(t *testing.T) {
	body, err := EncodeJSONBody(map[string]string{"name": "demo"})
	if err != nil {
		t.Fatalf("EncodeJSONBody() error = %v", err)
	}
	if got := string(body); !strings.Contains(got, `"name":"demo"`) {
		t.Fatalf("body = %q, want JSON containing name=demo", got)
	}
}

func TestEncodeJSONBody_Error(t *testing.T) {
	type bad struct {
		Fn func()
	}
	_, err := EncodeJSONBody(bad{Fn: func() {}})
	if err == nil {
		t.Fatal("EncodeJSONBody() error = nil, want non-nil")
	}
}
