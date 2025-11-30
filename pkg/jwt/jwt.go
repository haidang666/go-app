package jwt

import (
	"errors"
	"time"

	jwtV5 "github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Client struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTClient(secretKey string, tokenDuration time.Duration) *Client {
	return &Client{
		secretKey,
		tokenDuration,
	}
}

func (c *Client) Generate(claims jwtV5.Claims) (string, error) {
	token := jwtV5.NewWithClaims(jwtV5.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(c.secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (c *Client) Verify(tokenStr string, claims jwtV5.Claims) error {
	token, err := jwtV5.ParseWithClaims(tokenStr, claims, func(t *jwtV5.Token) (any, error) {
		if _, ok := t.Method.(*jwtV5.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(c.secretKey), nil
	})
	if err != nil || !token.Valid {
		return ErrInvalidToken
	}
	return nil
}
