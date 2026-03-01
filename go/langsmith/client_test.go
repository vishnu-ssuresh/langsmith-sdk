package langsmith

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewClient_ReturnsErrInvalidConfigWhenAPIKeyMissing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LANGSMITH_API_KEY", "")
	t.Setenv("LANGSMITH_WORKSPACE_ID", "")
	t.Setenv("LANGSMITH_ENDPOINT", "")

	_, err := NewClient(ClientOptions{})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("NewClient() error = %v, want ErrInvalidConfig", err)
	}
}

func TestNewClient_UsesEnvAPIKey(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LANGSMITH_API_KEY", "env-key")
	t.Setenv("LANGSMITH_WORKSPACE_ID", "")
	t.Setenv("LANGSMITH_ENDPOINT", "")

	client, err := NewClient(ClientOptions{})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClient_UsesConfigAPIKeyWhenEnvMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("LANGSMITH_API_KEY", "")
	t.Setenv("LANGSMITH_WORKSPACE_ID", "")
	t.Setenv("LANGSMITH_ENDPOINT", "")

	configDir := filepath.Join(home, ".langsmith-cli")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("api-key: config-key\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	client, err := NewClient(ClientOptions{})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}
