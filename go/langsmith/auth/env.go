package auth

import (
	"context"
	"os"
	"strings"
)

// EnvResolver reads credentials from process environment variables.
type EnvResolver struct {
	lookup func(string) (string, bool)
}

// NewEnvResolver creates a resolver backed by os.LookupEnv.
func NewEnvResolver() *EnvResolver {
	return &EnvResolver{lookup: os.LookupEnv}
}

func newEnvResolverWithLookup(lookup func(string) (string, bool)) *EnvResolver {
	return &EnvResolver{lookup: lookup}
}

// Resolve returns credentials from LANGSMITH_* environment variables.
func (r *EnvResolver) Resolve(_ context.Context) (Credentials, error) {
	return Credentials{
		APIKey:      firstLookup(r.lookup, "LANGSMITH_API_KEY"),
		WorkspaceID: firstLookup(r.lookup, "LANGSMITH_WORKSPACE_ID"),
		Endpoint:    firstLookup(r.lookup, "LANGSMITH_ENDPOINT"),
	}, nil
}

func firstLookup(lookup func(string) (string, bool), keys ...string) string {
	for _, key := range keys {
		value, ok := lookup(key)
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

func normalizeValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return strings.TrimSpace(value)
}
