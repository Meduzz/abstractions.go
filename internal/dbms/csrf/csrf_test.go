package csrf

import (
	"context"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDbmsCsrf(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}

	subject, err := NewDbmsCSRFStorageDelegate(db, time.Second)

	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	token := &lib.CSRFToken{
		Key:   "key",
		Value: "value",
	}

	t.Run("happy chappy", func(t *testing.T) {
		err := subject.Store(ctx, token)

		if err != nil {
			t.Error("generating token threw error", err)
		}

		valid, err := subject.Verify(ctx, token)

		if err != nil {
			t.Error("verifying token threw error", err)
		}

		if !valid {
			t.Error("token was not valid")
		}

		valid, err = subject.Verify(ctx, token)

		if err != nil {
			t.Error("revalidaton threw error", err)
		}

		if valid {
			t.Error("the token was still valid")
		}
	})

	t.Run("verify garbage", func(t *testing.T) {
		valid, err := subject.Verify(ctx, &lib.CSRFToken{})

		if err != nil {
			t.Error("invalid data threw error", err)
		}

		if valid {
			t.Error("invalid data was valid")
		}
	})

	t.Run("slow verifier is slow", func(t *testing.T) {
		err := subject.Store(ctx, token)

		if err != nil {
			t.Error("generating token threw error", err)
		}

		<-time.After(time.Second)

		valid, err := subject.Verify(ctx, token)

		if err != nil {
			t.Error("verifying token threw error", err)
		}

		if valid {
			t.Error("timed out token was valid")
		}
	})
}
