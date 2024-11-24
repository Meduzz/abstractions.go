package csrf

import (
	"context"
	"time"

	"github.com/Meduzz/abstractions.go/internal/interval"
	"github.com/Meduzz/abstractions.go/lib"
)

type (
	localToken struct {
		expires time.Time
		token   *lib.CSRFToken
	}

	localCsrf struct {
		storage map[string]*localToken
		ttl     time.Duration
	}
)

func NewLocalCsrf(ttl time.Duration) lib.CSRFStorageDelegate {
	storage := make(map[string]*localToken)

	interval.OnInterval(5*time.Second, func() {
		for k, v := range storage {
			if v.expires.Before(time.Now()) {
				delete(storage, k)
			}
		}
	})

	return &localCsrf{storage, ttl}
}

func (l *localCsrf) Store(ctx context.Context, token *lib.CSRFToken) error {
	storedToken := &localToken{
		token:   token,
		expires: time.Now().Add(l.ttl),
	}

	l.storage[token.Key] = storedToken

	return nil
}

func (l *localCsrf) Verify(ctx context.Context, token *lib.CSRFToken) (bool, error) {
	storedToken, ok := l.storage[token.Key]

	if !ok {
		return false, nil
	}

	if storedToken.expires.Before(time.Now()) {
		return false, nil
	}

	delete(l.storage, token.Key)

	return true, nil
}
