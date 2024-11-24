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

	// CSRFStorageDelegate - to deal with storing and verifying existance of tokens
	CSRFStorageDelegate interface {
		Store(context.Context, *CSRFToken) error
		Verify(context.Context, *CSRFToken) (bool, error)
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
