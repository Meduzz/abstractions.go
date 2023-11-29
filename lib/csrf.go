package lib

import (
	"context"
	"time"
)

type (
	CSRFToken struct {
		Key   string
		Value string
	}

	// CSRFAbstraction - to create and verify one time CSRF tokens.
	CSRFAbstraction interface {
		// Generate - generates a new random CSRF token.
		Generate(ctx context.Context, duration time.Duration) (*CSRFToken, error)
		// Verify - verifies a CSRF token, returns false if provided values does not match.
		Verify(ctx context.Context, key, value string) (bool, error)
	}
)
