package auth

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultConfigPath = "~/.langsmith-cli/config.yaml"

// ConfigResolver reads credentials from a local YAML config file.
type ConfigResolver struct {
	path     string
	readFile func(string) ([]byte, error)
}

// NewConfigResolver creates a resolver for the provided config path.
//
// If path is empty, ~/.langsmith-cli/config.yaml is used.
func NewConfigResolver(path string) *ConfigResolver {
	return &ConfigResolver{
		path:     path,
		readFile: os.ReadFile,
	}
}

func newConfigResolverWithReader(path string, readFile func(string) ([]byte, error)) *ConfigResolver {
	return &ConfigResolver{
		path:     path,
		readFile: readFile,
	}
}

// Resolve returns credentials loaded from config.
func (r *ConfigResolver) Resolve(_ context.Context) (Credentials, error) {
	path, err := resolveConfigPath(r.path)
	if err != nil {
		return Credentials{}, err
	}

	data, err := r.readFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Credentials{}, nil
		}
		return Credentials{}, fmt.Errorf("auth: read config: %w", err)
	}

	values := parseFlatYAML(data)
	return Credentials{
		APIKey: firstMapValue(values,
			"api-key",
			"api_key",
			"langsmith-api-key",
			"langsmith_api_key",
		),
		WorkspaceID: firstMapValue(values,
			"workspace-id",
			"workspace_id",
			"langsmith-workspace-id",
			"langsmith_workspace_id",
		),
		Endpoint: firstMapValue(values,
			"endpoint",
			"base-url",
			"base_url",
			"langsmith-endpoint",
			"langsmith_endpoint",
		),
	}, nil
}

func resolveConfigPath(path string) (string, error) {
	if path == "" {
		path = defaultConfigPath
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("auth: resolve home dir: %w", err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path, nil
}

func parseFlatYAML(data []byte) map[string]string {
	out := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexRune(line, ':')
		if idx <= 0 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(line[:idx]))
		value := normalizeValue(line[idx+1:])
		if key == "" || value == "" {
			continue
		}
		out[key] = value
	}
	return out
}

func firstMapValue(values map[string]string, keys ...string) string {
	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			continue
		}
		value = normalizeValue(value)
		if value != "" {
			return value
		}
	}
	return ""
}
