package abstractions

import (
	"time"

	"github.com/Meduzz/abstractions.go/internal/local/caching"
	"github.com/Meduzz/abstractions.go/internal/local/csrf"
	"github.com/Meduzz/abstractions.go/internal/local/eventing"
	"github.com/Meduzz/abstractions.go/internal/local/log"
	"github.com/Meduzz/abstractions.go/lib"
)

// RedisCaching - create a new caching module with the provided codec.
func LocalCaching[T any](ttl time.Duration) lib.CacheAbstraction[T] {
	return caching.NewCache[T](ttl)
}

// LocalCSRF - create a new CSRF module.
func LocalCSRF(ttl time.Duration) lib.CSRFAbstraction {
	return csrf.NewLocalCsrf(ttl)
}

// LocalEventing - create a new eventing module.
func Localventing[T any]() lib.EventingAbstraction[T] {
	transport := make(chan *T)
	return eventing.NewLocalEventing(transport)
}

// LocalLog - create a local work log.
func LocalLog[T any]() lib.LogAbstraction[T] {
	return log.NewLocalLog[T]()
}
