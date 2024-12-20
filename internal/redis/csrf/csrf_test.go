package csrf_test

import (
	"context"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/internal/redis/csrf"
	"github.com/Meduzz/abstractions.go/lib"
	"github.com/Meduzz/abstractions.go/lib/specific"
	"github.com/Meduzz/helper/rudis"
)

func TestCSRF(t *testing.T) {
	conn := rudis.Connect()
	cfg := specific.NewRedisConfig(conn, "")
	ctx := context.Background()
	module := csrf.NewRedisCSRFStorageDelegate(cfg, time.Second, "csrf")
	token := &lib.CSRFToken{
		Key:   "key",
		Value: "value",
	}

	defer conn.Close()

	t.Run("generate and verify a token", func(t *testing.T) {
		err := module.Store(ctx, token)

		if err != nil {
			t.Error("generating token threw error", err)
		}

		valid, err := module.Verify(ctx, token)

		if err != nil {
			t.Error("verifying token threw error", err)
		}

		if !valid {
			t.Error("token was not valid")
		}

		valid, err = module.Verify(ctx, token)

		if err != nil {
			t.Error("revalidaton threw error", err)
		}

		if valid {
			t.Error("the token was still valid")
		}
	})

	t.Run("verify garbage", func(t *testing.T) {
		valid, err := module.Verify(ctx, &lib.CSRFToken{})

		if err != nil {
			t.Error("invalid data threw error", err)
		}

		if valid {
			t.Error("invalid data was valid")
		}
	})

	t.Run("slow verifier is slow", func(t *testing.T) {
		err := module.Store(ctx, token)

		if err != nil {
			t.Error("generating token threw error", err)
		}

		<-time.After(time.Second)

		valid, err := module.Verify(ctx, token)

		if err != nil {
			t.Error("verifying token threw error", err)
		}

		if valid {
			t.Error("timed out token was valid")
		}
	})
}
