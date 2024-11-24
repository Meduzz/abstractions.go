package abstractions

import (
	"time"

	"github.com/Meduzz/abstractions.go/internal/local/caching"
	"github.com/Meduzz/abstractions.go/internal/local/csrf"
	"github.com/Meduzz/abstractions.go/lib"
)

// LocalCachingDelegate - create a new caching delegate with the provided settings.
func LocalCachingDelegate(eviction lib.Eviction, ttl time.Duration) lib.CacheStorageDelegate {
	return caching.NewCache(eviction, ttl)
}

// LocalCSRFDelegate - create a new CSRF storage delegate.
func LocalCSRFDelegate(ttl time.Duration) lib.CSRFStorageDelegate {
	return csrf.NewLocalCsrf(ttl)
}
