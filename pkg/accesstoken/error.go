package accesstoken

import "errors"

var (
	ErrExpiredToken    = errors.New("expired token")
	ErrInvalidAudience = errors.New("invalid audience")
	ErrInvalidIssuer   = errors.New("invalid issuer")
)
