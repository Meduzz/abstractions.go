package abstractions

import (
	"time"

	root "github.com/Meduzz/abstractions.go/internal/redis"
	"github.com/Meduzz/abstractions.go/internal/redis/caching"
	"github.com/Meduzz/abstractions.go/internal/redis/csrf"
	"github.com/Meduzz/abstractions.go/internal/redis/eventing"
	"github.com/Meduzz/abstractions.go/internal/redis/log"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/go-redis/redis/v8"
)

// CreateRedisConfig - create a RedisConfig from a redis connection and a prefix (optionally empty).
func CreateRedisConfig(conn *redis.Client, prefix string) *root.RedisConfig {
	return root.NewRedisConfig(conn, prefix)
}

// RedisCaching - create a new caching module with the provided config and codec.
func RedisCaching[T any](config *root.RedisConfig, codec lib.Codec[T], ttl time.Duration, name string) lib.CacheAbstraction[T] {
	return caching.NewCaching[T](config, codec, ttl, name)
}

// RedisCSRF - create a new CSRF module with the provided config.
func RedisCSRF(config *root.RedisConfig, ttl time.Duration) lib.CSRFAbstraction {
	return csrf.NewCSRFAbstraction(config, ttl)
}

// RedisEventing - create a new eventing module.
func RedisEventing[T any](topic string, codec lib.Codec[T], config *root.RedisConfig) lib.EventingAbstraction[T] {
	return eventing.NewEventing[T](topic, codec, config)
}

// RedisLog - create a new redis based work log.
func RedisLog[T any](config *root.RedisConfig, codec lib.Codec[T], name string) lib.LogAbstraction[T] {
	return log.NewRedisLog(config, codec, name)
}
