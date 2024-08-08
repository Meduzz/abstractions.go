package csrf_test

import (
	"context"
	"testing"
	"time"

	root "github.com/Meduzz/abstractions.go/internal/redis"
	"github.com/Meduzz/abstractions.go/internal/redis/csrf"
	"github.com/Meduzz/helper/rudis"
)

func TestCSRF(t *testing.T) {
	conn := rudis.Connect()
	cfg := root.NewRedisConfig(conn, "")
	ctx := context.Background()
	module := csrf.NewCSRFAbstraction(cfg)

	defer conn.Close()

	t.Run("generate and verify a token", func(t *testing.T) {
		token, err := module.Generate(ctx, time.Second)

		if err != nil {
			t.Error("generating token threw error", err)
		}

		if token == nil {
			t.Error("token was nil")
			return
		}

		valid, err := module.Verify(ctx, token.Key, token.Value)

		if err != nil {
			t.Error("verifying token threw error", err)
		}

		if !valid {
			t.Error("token was not valid")
		}

		valid, err = module.Verify(ctx, token.Key, token.Value)

		if err != nil {
			t.Error("revalidaton threw error", err)
		}

		if valid {
			t.Error("the token was still valid")
		}
	})

	t.Run("verify garbage", func(t *testing.T) {
		valid, err := module.Verify(ctx, "bad", "data")

		if err != nil {
			t.Error("invalid data threw error", err)
		}

		if valid {
			t.Error("invalid data was valid")
		}
	})

	t.Run("slow verifier is slow", func(t *testing.T) {
		token, err := module.Generate(ctx, time.Second)

		if err != nil {
			t.Error("generating token threw error", err)
		}

		<-time.After(time.Second)

		valid, err := module.Verify(ctx, token.Key, token.Value)

		if err != nil {
			t.Error("verifying token threw error", err)
		}

		if valid {
			t.Error("timed out token was valid")
		}
	})
}
