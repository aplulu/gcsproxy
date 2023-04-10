package accesstoken

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Audience []string

func (aud *Audience) UnmarshalJSON(data []byte) error {
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch v := value.(type) {
	case []string:
		*aud = v
	case string:
		*aud = append(*aud, v)
	case []interface{}:
		for _, vv := range v {
			vs, ok := vv.(string)
			if !ok {
				return &json.UnsupportedTypeError{Type: reflect.TypeOf(vv)}
			}
			*aud = append(*aud, vs)
		}
	case nil:
		return nil
	default:
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(v)}
	}

	return nil
}

func (aud *Audience) MarshalJSON() ([]byte, error) {
	if len(*aud) == 1 {
		return json.Marshal((*aud)[0])
	}
	return json.Marshal([]string(*aud))
}

type AccessToken struct {
	Issuer         string   `json:"iss"`
	ExpirationTime int64    `json:"exp"`
	Audience       Audience `json:"aud"`
	Subject        string   `json:"sub"`
	IssuedAt       int64    `json:"iat"`
}

func (a *AccessToken) Valid() error {
	now := time.Now().UTC()

	exp := time.Unix(a.ExpirationTime, 0)
	if !now.Before(exp) {
		return fmt.Errorf("accesstoken.Valid: expired access token: %w", ErrExpiredToken)
	}

	return nil
}

func (a *AccessToken) Sign(secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	ts, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("accesstoken.Sign: failed to sign token: %w", err)
	}

	return ts, nil
}

func ParseAccessToken(accessToken string, issuer string, audience string, secretKey []byte) (*AccessToken, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("accesstoken.ParseAccessToken: unknown signing method: %v", token.Method)
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("accesstoken.ParseAccessToken: failed to parse token: %w", err)
	}
	claims, ok := token.Claims.(*AccessToken)
	if !ok {
		return nil, fmt.Errorf("accesstoken.ParseAccessToken: failed to convert access token")
	}

	var found bool
	for _, v := range claims.Audience {
		if subtle.ConstantTimeCompare([]byte(v), []byte(audience)) == 1 {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("accesstoken.ParseAccessToken: invalid audience: %w", ErrInvalidAudience)
	}

	if subtle.ConstantTimeCompare([]byte(claims.Issuer), []byte(issuer)) == 0 {
		return nil, ErrInvalidIssuer
	}

	return claims, err
}
