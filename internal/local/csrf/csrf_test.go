package csrf_test

import (
	"context"
	"testing"
	"time"

	"github.com/Meduzz/abstractions.go/internal/local/csrf"
	"github.com/Meduzz/abstractions.go/lib"
)

func TestLocalCsrf(t *testing.T) {
	subject := csrf.NewLocalCsrf(time.Second)
	ctx := context.Background()
	token := &lib.CSRFToken{
		Key:   "key",
		Value: "value",
	}

	t.Run("happy campers", func(t *testing.T) {
		subject.Store(ctx, token)
		ok, err := subject.Verify(ctx, token)

		if err != nil {
			t.Errorf("error was not nil but %v", err)
		}

		if !ok {
			t.Error("the token was not valid")
		}

		// verify the token a second time should fail
		ok, err = subject.Verify(ctx, token)

		if err == nil {
			t.Error("there was no error")
		}

		if ok {
			t.Error("the token was valid")
		}
	})

	t.Run("verify garbage", func(t *testing.T) {
		bad := &lib.CSRFToken{
			Key:   "asdf",
			Value: "qwer",
		}
		ok, err := subject.Verify(ctx, bad)

		if err == nil {
			t.Error("there was no error")
		}

		if ok {
			t.Error("the token was valid")
		}
	})

}
