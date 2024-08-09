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

	t.Run("happy campers", func(t *testing.T) {
		token, _ := subject.Generate(ctx)

		ok, err := subject.Verify(ctx, token.Key, token.Value)

		if err != nil {
			t.Errorf("error was not nil but %v", err)
		}

		if !ok {
			t.Error("the token was not valid")
		}

		// verify the token a second time should fail
		ok, err = subject.Verify(ctx, token.Key, token.Value)

		if err == nil {
			t.Error("there was no error")
		}

		if ok {
			t.Error("the token was valid")
		}
	})

	t.Run("verify garbage", func(t *testing.T) {
		ok, err := subject.Verify(ctx, "asdf", "qwer")

		if err == nil {
			t.Error("there was no error")
		}

		if ok {
			t.Error("the token was valid")
		}
	})

	t.Run("playing with fire (jwt)", func(t *testing.T) {
		token, _ := subject.Generate(ctx)
		strToken, err := token.Encode()

		if err != nil {
			t.Errorf("generating jwt threw error: %v", err)
		}

		result, err := lib.DecodeToken(strToken)

		if err != nil {
			t.Errorf("parsing token threw error: %v", err)
		}

		if result == nil {
			t.Error("token was not parsed")
		}

		if result.Key != token.Key || result.Value != token.Value {
			t.Error("token and result does not match")
		}
	})

}
