package lib

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type (
	CSRFToken struct {
		Key   string `json:"key"`
		Value string `json:"form"`
	}

	// CSRFAbstraction - to create and verify one time CSRF tokens.
	CSRFAbstraction interface {
		// Generate - generates a new random CSRF token.
		Generate(ctx context.Context) (*CSRFToken, error)
		// Verify - verifies a CSRF token, returns false if provided values does not match.
		Verify(ctx context.Context, key, value string) (bool, error)
	}
)

func (c *CSRFToken) Encode() (string, error) {
	data := jwt.MapClaims{}
	data["key"] = c.Key
	data["value"] = c.Value
	token := jwt.NewWithClaims(jwt.SigningMethodNone, data)

	return token.SignedString(jwt.UnsafeAllowNoneSignatureType)
}

func DecodeToken(token string) (*CSRFToken, error) {
	data := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, data, func(t *jwt.Token) (interface{}, error) {
		return jwt.UnsafeAllowNoneSignatureType, nil
	}, jwt.WithValidMethods([]string{"none"}))

	if err != nil {
		return nil, err
	}

	result := &CSRFToken{}
	result.Key = data["key"].(string)
	result.Value = data["value"].(string)

	return result, nil
}
