package abstractions

import (
	"time"

	"github.com/Meduzz/abstractions.go/internal/redis/caching"
	"github.com/Meduzz/abstractions.go/internal/redis/csrf"
	"github.com/Meduzz/abstractions.go/internal/redis/log"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/go-redis/redis/v8"
)

// CreateRedisConfig - create a RedisConfig from a redis connection and a prefix (optionally empty).
func CreateRedisConfig(conn *redis.Client, prefix string) *specific.RedisConfig {
	return specific.NewRedisConfig(conn, prefix)
}

// RedisCachingDelegate - create a new caching module with the provided settings.
func RedisCachingDelegate(config *specific.RedisConfig, eviction lib.Eviction, ttl time.Duration, name string) lib.CacheStorageDelegate {
	return caching.NewCaching(config, eviction, ttl, name)
}

// RedisCSRFDelegate - create a new CSRF storage delegate with the provided settings.
func RedisCSRFDelegate(config *specific.RedisConfig, ttl time.Duration, name string) lib.CSRFStorageDelegate {
	return csrf.NewRedisCSRFStorageDelegate(config, ttl, name)
}

// RedisWorkLogDelegate - create a new redis based work log delegate.
func RedisWorkLogDelegate(config *specific.RedisConfig, name string) lib.WorkLogDelegate {
	return log.NewRedisWorkLog(config, name)
}
