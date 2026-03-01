package auth

import (
	"context"
	"errors"
	"testing"
)

type fakeResolver struct {
	creds Credentials
	err   error
}

func (r fakeResolver) Resolve(context.Context) (Credentials, error) {
	if r.err != nil {
		return Credentials{}, r.err
	}
	return r.creds, nil
}

func TestResolveCredentials_MergesByPrecedence(t *testing.T) {
	creds, err := ResolveCredentials(
		context.Background(),
		fakeResolver{creds: Credentials{Endpoint: "https://explicit.example.com"}},
		fakeResolver{creds: Credentials{APIKey: "env-key", WorkspaceID: "env-workspace", Endpoint: "https://env.example.com"}},
		fakeResolver{creds: Credentials{APIKey: "config-key", WorkspaceID: "config-workspace", Endpoint: "https://config.example.com"}},
	)
	if err != nil {
		t.Fatalf("ResolveCredentials() error = %v", err)
	}

	if creds.APIKey != "env-key" {
		t.Fatalf("APIKey = %q, want %q", creds.APIKey, "env-key")
	}
	if creds.WorkspaceID != "env-workspace" {
		t.Fatalf("WorkspaceID = %q, want %q", creds.WorkspaceID, "env-workspace")
	}
	if creds.Endpoint != "https://explicit.example.com" {
		t.Fatalf("Endpoint = %q, want %q", creds.Endpoint, "https://explicit.example.com")
	}
}

func TestResolveCredentials_RequiresAPIKey(t *testing.T) {
	_, err := ResolveCredentials(
		context.Background(),
		fakeResolver{creds: Credentials{WorkspaceID: "workspace"}},
	)
	if !errors.Is(err, ErrCredentialsNotFound) {
		t.Fatalf("ResolveCredentials() error = %v, want ErrCredentialsNotFound", err)
	}
}

func TestResolveCredentials_ResolverError(t *testing.T) {
	wantErr := errors.New("failed resolver")
	_, err := ResolveCredentials(
		context.Background(),
		fakeResolver{err: wantErr},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("ResolveCredentials() error = %v, want %v", err, wantErr)
	}
}
