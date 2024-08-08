package csrf

import (
	"context"
	"time"

	"github.com/Meduzz/abstractions.go/internal/local/interval"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/helper/hashing"
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

func NewLocalCsrf(ttl time.Duration) lib.CSRFAbstraction {
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

func (l *localCsrf) Generate(ctx context.Context) (*lib.CSRFToken, error) {
	realToken := &lib.CSRFToken{
		Key:   hashing.Token(),
		Value: hashing.Secret(),
	}
	storedToken := &localToken{
		token:   realToken,
		expires: time.Now().Add(l.ttl),
	}

	l.storage[realToken.Key] = storedToken

	return realToken, nil
}

func (l *localCsrf) Verify(ctx context.Context, key, value string) (bool, error) {
	storedToken, ok := l.storage[key]

	if !ok {
		return false, lib.ErrKeyNotFound
	}

	if storedToken.expires.Before(time.Now()) {
		return false, nil
	}

	delete(l.storage, key)

	return true, nil
}
