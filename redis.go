package abstractions

import (
	"github.com/Meduzz/abstractions.go/internal/redis/caching"
	"github.com/Meduzz/abstractions.go/internal/redis/csrf"
	"github.com/Meduzz/abstractions.go/internal/redis/eventing"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/go-redis/redis/v8"
)

// CreateConfig - create a RedisConfig from a redis connection and a prefix (optionally empty).
func CreateConfig(conn *redis.Client, prefix string) *lib.RedisConfig {
	return lib.NewRedisConfig(conn, prefix)
}

// Caching - create a new caching module with the provided config and codec.
func Caching[T any](config *lib.RedisConfig, codec lib.Codec[T]) lib.CacheAbstraction[T] {
	return caching.NewCaching[T](config, codec)
}

// Csrf - create a new CSRF module with the provided config.
func Csrf(config *lib.RedisConfig) lib.CSRFAbstraction {
	return csrf.NewCSRFAbstraction(config)
}

// Eventing - create a new eventing module.
func Eventing(config *lib.RedisConfig) lib.EventingAbstraction {
	return eventing.NewEventing(config)
}
