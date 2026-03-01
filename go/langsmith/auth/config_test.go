package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigResolver_Resolve(t *testing.T) {
	readFile := func(path string) ([]byte, error) {
		if path != "/tmp/config.yaml" {
			t.Fatalf("path = %q, want %q", path, "/tmp/config.yaml")
		}
		return []byte(`
api-key: test-api-key
workspace_id: workspace-1
base_url: https://api.smith.langchain.com
`), nil
	}

	r := newConfigResolverWithReader("/tmp/config.yaml", readFile)
	creds, err := r.Resolve(context.Background())
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if creds.APIKey != "test-api-key" {
		t.Fatalf("APIKey = %q, want %q", creds.APIKey, "test-api-key")
	}
	if creds.WorkspaceID != "workspace-1" {
		t.Fatalf("WorkspaceID = %q, want %q", creds.WorkspaceID, "workspace-1")
	}
	if creds.Endpoint != "https://api.smith.langchain.com" {
		t.Fatalf("Endpoint = %q, want %q", creds.Endpoint, "https://api.smith.langchain.com")
	}
}

func TestConfigResolver_Resolve_MissingFile(t *testing.T) {
	r := newConfigResolverWithReader("/tmp/missing.yaml", func(string) ([]byte, error) {
		return nil, os.ErrNotExist
	})

	creds, err := r.Resolve(context.Background())
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if creds.APIKey != "" || creds.WorkspaceID != "" || creds.Endpoint != "" {
		t.Fatalf("creds = %+v, want empty", creds)
	}
}

func TestConfigResolver_Resolve_ReadError(t *testing.T) {
	r := newConfigResolverWithReader("/tmp/bad.yaml", func(string) ([]byte, error) {
		return nil, errors.New("boom")
	})

	_, err := r.Resolve(context.Background())
	if err == nil {
		t.Fatal("Resolve() error = nil, want non-nil")
	}
}

func TestResolveConfigPath_DefaultAndTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := resolveConfigPath("")
	if err != nil {
		t.Fatalf("resolveConfigPath(\"\") error = %v", err)
	}
	want := filepath.Join(home, ".langsmith-cli", "config.yaml")
	if got != want {
		t.Fatalf("resolveConfigPath(\"\") = %q, want %q", got, want)
	}

	got, err = resolveConfigPath("~/custom/path.yaml")
	if err != nil {
		t.Fatalf("resolveConfigPath(\"~/custom/path.yaml\") error = %v", err)
	}
	want = filepath.Join(home, "custom", "path.yaml")
	if got != want {
		t.Fatalf("resolveConfigPath(\"~/custom/path.yaml\") = %q, want %q", got, want)
	}
}
