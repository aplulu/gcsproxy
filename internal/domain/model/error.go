package model

import "errors"

var (
	ErrInvalidIDToken      = errors.New("invalid id token")
	ErrInvalidRedirectURL  = errors.New("invalid redirect URL")
	ErrInvalidState        = errors.New("invalid state")
	ErrInvalidHostedDomain = errors.New("invalid hosted domain")
)
