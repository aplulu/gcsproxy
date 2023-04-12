package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"

	"github.com/aplulu/gcsproxy/internal/config"
	"github.com/aplulu/gcsproxy/internal/util"
	"github.com/aplulu/gcsproxy/pkg/accesstoken"
)

type IDTokenClaims struct {
	Sub string `json:"sub"`
}

var oidcProvider *oidc.Provider
var oidcVerifier *oidc.IDTokenVerifier

// GetOIDCConfig returns OIDC config
func GetOIDCConfig(ctx context.Context) (*oauth2.Config, error) {
	if oidcVerifier == nil {
		var err error
		oidcProvider, err = oidc.NewProvider(ctx, config.OIDCProvider())
		if err != nil {
			return nil, err
		}
		oidcVerifier = oidcProvider.Verifier(&oidc.Config{
			ClientID: config.OIDCClientID(),
		})
	}

	var endpoint oauth2.Endpoint
	if config.OIDCAuthorizeURL() != "" && config.OIDCTokenURL() != "" {
		endpoint = oauth2.Endpoint{
			AuthURL:  config.OIDCAuthorizeURL(),
			TokenURL: config.OIDCTokenURL(),
		}
	} else {
		endpoint = oidcProvider.Endpoint()
	}

	return &oauth2.Config{
		ClientID:     config.OIDCClientID(),
		ClientSecret: config.OIDCClientSecret(),
		Endpoint:     endpoint,
		RedirectURL:  config.BaseURL() + "/_gcsproxy/oidc/callback",
		Scopes:       config.OIDCScopes(),
	}, nil
}

// ExchangeOIDCToken creates token for OIDC Authenticate session
func ExchangeOIDCToken(ctx context.Context, code string) (*IDTokenClaims, error) {
	oc, err := GetOIDCConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("model.ExchangeOIDCToken: failed to retrive OAuth2 config: %w", err)
	}

	// アクセストークン取得
	token, err := oc.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("model.ExchangeOIDCToken: failed to exchange token: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("model.ExchangeOIDCToken: missing id token: %w", ErrInvalidIDToken)
	}

	idToken, err := oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("model.ExchangeOIDCToken: failed to verify IDToken: %w", ErrInvalidIDToken)
	}

	claims := new(IDTokenClaims)
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("model.ExchangeOIDCToken: failed to parse claims: %w", err)
	}

	return claims, nil
}

type OIDCSession struct {
	State       string
	RedirectURL string
}

func (s *OIDCSession) Valid() error {
	return nil
}

// NewOIDCSession creates new OIDCSession and returns JWT token.
func NewOIDCSession(redirectURL string) (string, *OIDCSession, error) {
	if !strings.HasPrefix(redirectURL, "/") {
		return "", nil, fmt.Errorf("model.NewOIDCSession: redirectURL must start with /: %s: %w", redirectURL, ErrInvalidRedirectURL)
	}

	sess := &OIDCSession{
		State:       util.RandomString(32),
		RedirectURL: redirectURL,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sess)
	signedToken, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		return "", nil, fmt.Errorf("model.NewOIDCSession: failed to sign session payload: %w", err)
	}
	return signedToken, sess, nil
}

// ParseOIDCSession parses OIDCSession from JWT token.
func ParseOIDCSession(tokenString string) (*OIDCSession, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OIDCSession{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("model.ParseOIDCSession: unknown signing method: %v", token.Method)
		}

		return []byte(config.JWTSecret()), nil
	})
	if err != nil {
		return nil, fmt.Errorf("model.ParseOIDCSession: failed to parse jwt: %w", err)
	}

	claims, ok := token.Claims.(*OIDCSession)
	if !ok {
		return nil, fmt.Errorf("model.ParseOIDCSession: failed to convert session token")
	}

	return claims, nil
}

// CreateAuthSession creates JWT token for authentication.
func CreateAuthSession(id string) (string, *time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Duration(config.JWTExpiration()) * time.Second)
	at := &accesstoken.AccessToken{
		Issuer:         config.BaseURL(),
		ExpirationTime: exp.Unix(),
		Audience:       []string{config.BaseURL()},
		Subject:        id,
		IssuedAt:       now.Unix(),
	}

	token, err := at.Sign([]byte(config.JWTSecret()))
	if err != nil {
		return "", nil, fmt.Errorf("model.CreateAuthSession: failed to sign token: %w", err)
	}

	return token, &exp, nil
}
