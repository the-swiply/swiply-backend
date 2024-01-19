package auf

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	secretInstance []byte
)

func SetSecret(secret []byte) {
	secretInstance = secret
}

type JWTAccessProperties struct {
	User        string
	TTL         time.Duration
	Fingerprint string
}

type JWTRefreshProperties struct {
	User string
	TTL  time.Duration
}

func GenerateAccessJWT(props JWTAccessProperties) (string, error) {
	return GenerateAccessJWTWithCustomSecret(props, secretInstance)
}

func GenerateRefreshJWT(props JWTRefreshProperties) (string, error) {
	return GenerateRefreshJWTWithCustomSecret(props, secretInstance)
}

func ValidateJWTAndExtractClaims(token string) (map[string]any, error) {
	return ValidateJWTAndExtractClaimsWithCustomSecret(token, secretInstance)
}

func GenerateAccessJWTWithCustomSecret(props JWTAccessProperties, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":        props.User,
		"fingerprint": props.Fingerprint,
		"exp":         time.Now().Add(props.TTL).Unix(),
	})

	return token.SignedString(secret)
}

func GenerateRefreshJWTWithCustomSecret(props JWTRefreshProperties, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": props.User,
		"exp":  time.Now().Add(props.TTL).Unix(),
	})

	return token.SignedString(secret)
}

func ValidateJWTAndExtractClaimsWithCustomSecret(token string, secret []byte) (map[string]any, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return secret, nil
	})

	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, errors.New("invalid token format")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, errors.New("invalid token signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, errors.New("token has expired")
	case t != nil && t.Valid:
		claims, ok := t.Claims.(jwt.MapClaims)
		if ok {
			return claims, nil
		} else {
			return nil, errors.New("can't extract claims")
		}
	default:
		return nil, fmt.Errorf("can't validate token: %w", err)
	}
}
