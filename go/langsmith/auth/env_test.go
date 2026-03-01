package auth

import (
	"context"
	"testing"
)

func TestEnvResolver_Resolve(t *testing.T) {
	lookup := func(key string) (string, bool) {
		switch key {
		case "LANGSMITH_API_KEY":
			return " api-key ", true
		case "LANGSMITH_WORKSPACE_ID":
			return "\"workspace-id\"", true
		case "LANGSMITH_ENDPOINT":
			return "https://api.example.com", true
		default:
			return "", false
		}
	}

	r := newEnvResolverWithLookup(lookup)
	creds, err := r.Resolve(context.Background())
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if creds.APIKey != "api-key" {
		t.Fatalf("APIKey = %q, want %q", creds.APIKey, "api-key")
	}
	if creds.WorkspaceID != "workspace-id" {
		t.Fatalf("WorkspaceID = %q, want %q", creds.WorkspaceID, "workspace-id")
	}
	if creds.Endpoint != "https://api.example.com" {
		t.Fatalf("Endpoint = %q, want %q", creds.Endpoint, "https://api.example.com")
	}
}

func TestEnvResolver_Resolve_EmptyValues(t *testing.T) {
	lookup := func(key string) (string, bool) {
		switch key {
		case "LANGSMITH_API_KEY":
			return "   ", true
		case "LANGSMITH_WORKSPACE_ID":
			return "", false
		case "LANGSMITH_ENDPOINT":
			return "  ", true
		default:
			return "", false
		}
	}

	r := newEnvResolverWithLookup(lookup)
	creds, err := r.Resolve(context.Background())
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if creds.APIKey != "" || creds.WorkspaceID != "" || creds.Endpoint != "" {
		t.Fatalf("creds = %+v, want empty", creds)
	}
}
