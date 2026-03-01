package auth

import (
	"context"
	"errors"
)

// Resolver provides credentials for outgoing requests.
type Resolver interface {
	Resolve(ctx context.Context) (Credentials, error)
}

var (
	// ErrCredentialsNotFound means no resolver produced an API key.
	ErrCredentialsNotFound = errors.New("auth: credentials not found")
)

// ResolveCredentials merges credentials from resolvers in precedence order.
//
// Earlier resolvers take precedence over later ones for each individual field.
// At minimum, APIKey must be resolved for success.
func ResolveCredentials(ctx context.Context, resolvers ...Resolver) (Credentials, error) {
	var out Credentials
	for _, resolver := range resolvers {
		if resolver == nil {
			continue
		}
		creds, err := resolver.Resolve(ctx)
		if err != nil {
			return Credentials{}, err
		}
		mergeCredentials(&out, creds)
	}
	if out.APIKey == "" {
		return Credentials{}, ErrCredentialsNotFound
	}
	return out, nil
}

func mergeCredentials(dst *Credentials, src Credentials) {
	if dst.APIKey == "" && src.APIKey != "" {
		dst.APIKey = src.APIKey
	}
	if dst.WorkspaceID == "" && src.WorkspaceID != "" {
		dst.WorkspaceID = src.WorkspaceID
	}
	if dst.Endpoint == "" && src.Endpoint != "" {
		dst.Endpoint = src.Endpoint
	}
}
