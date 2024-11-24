package abstractions

import (
	"context"

	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/helper/hashing"
)

type (
	CSRF struct {
		storage lib.CSRFStorageDelegate
	}
)

func NewCSRF(storage lib.CSRFStorageDelegate) *CSRF {
	return &CSRF{storage}
}

func (c *CSRF) Generate(ctx context.Context) (*lib.CSRFToken, error) {
	key := hashing.Token()
	value := hashing.Secret()

	token := &lib.CSRFToken{
		Key:   key,
		Value: value,
	}

	err := c.storage.Store(ctx, token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *CSRF) Verify(ctx context.Context, token string) (bool, error) {
	csrfToken, err := lib.DecodeToken(token)

	if err != nil {
		return false, err
	}

	return c.storage.Verify(ctx, csrfToken)
}
