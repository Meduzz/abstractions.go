package specific

import (
	"fmt"

	"github.com/Meduzz/helper/fp/slice"
	"github.com/go-redis/redis/v8"
)

type (
	RedisConfig struct {
		conn   *redis.Client
		prefix string
	}
)

// NewRedisConfig - create a new config from the provided connection and prefix
func NewRedisConfig(conn *redis.Client, prefix string) *RedisConfig {
	return &RedisConfig{conn, prefix}
}

// Prefix - prefixes the provided key with the prefix of the config if set, or returns key.
func (c *RedisConfig) Prefix(keys ...string) string {
	if c.prefix != "" {
		return fmt.Sprintf("%s.%s", c.prefix, join(keys))
	} else {
		return join(keys)
	}
}

// Redis - return the connection of the config.
func (c *RedisConfig) Redis() *redis.Client {
	return c.conn
}

func join(strings []string) string {
	return slice.Fold(strings, "", func(in, agg string) string {
		if agg != "" {
			if in == "" {
				return agg
			}

			return fmt.Sprintf("%s.%s", agg, in)
		}

		return in
	})
}
