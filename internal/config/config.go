package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Listen                   string   `envconfig:"listen" default:""`
	Port                     string   `envconfig:"port" default:"8080"`
	GoogleCloudStorageBucket string   `envconfig:"google_cloud_storage_bucket"`
	MainPageSuffix           string   `envconfig:"main_page_suffix" default:"index.html"`
	NotFoundPage             string   `envconfig:"not_found_page" default:""`
	BaseURL                  string   `envconfig:"base_url" default:""`
	AuthType                 string   `envconfig:"auth_type" default:"none"`
	OIDCIssuer               string   `envconfig:"oidc_issuer" default:""`
	OIDCScopes               []string `envconfig:"oidc_scopes" default:"openid"`
	OIDCAuthorizeURL         string   `envconfig:"oidc_authorize_url" default:""`
	OIDCTokenURL             string   `envconfig:"oidc_token_url" default:""`
	OIDCClientID             string   `envconfig:"oidc_client_id" default:""`
	OIDCClientSecret         string   `envconfig:"oidc_client_secret" default:""`
	JWTExpiration            int64    `envconfig:"jwt_expiration" default:"3600"`
	JWTSecret                string   `envconfig:"jwt_secret"`
	BasicAuthUser            string   `envconfig:"basic_auth_user" default:""`
	BasicAuthPassword        string   `envconfig:"basic_auth_password" default:""`
}

var conf Config

// LoadConf Load Configurations
func LoadConf() error {
	if err := envconfig.Process("", &conf); err != nil {
		return fmt.Errorf("config.LoadConf: failed to load conf: %w", err)
	}

	return nil
}

func Listen() string {
	return conf.Listen
}

func Port() string {
	return conf.Port
}

func GoogleCloudStorageBucket() string {
	return conf.GoogleCloudStorageBucket
}

func MainPageSuffix() string {
	return conf.MainPageSuffix
}

func NotFoundPage() string {
	return conf.NotFoundPage
}

// BaseURL returns base URL
func BaseURL() string {
	return conf.BaseURL
}

func AuthType() string {
	return conf.AuthType
}

func OIDCIssuer() string {
	return conf.OIDCIssuer
}

func OIDCAuthorizeURL() string {
	return conf.OIDCAuthorizeURL
}

func OIDCTokenURL() string {
	return conf.OIDCTokenURL
}

func OIDCScopes() []string {
	return conf.OIDCScopes
}

func OIDCClientID() string {
	return conf.OIDCClientID
}

func OIDCClientSecret() string {
	return conf.OIDCClientSecret
}

func JWTExpiration() int64 {
	return conf.JWTExpiration
}

func JWTSecret() string {
	return conf.JWTSecret
}

func BasicAuthUser() string {
	return conf.BasicAuthUser
}

func BasicAuthPassword() string {
	return conf.BasicAuthPassword
}

func ValidateOIDC() error {
	if AuthType() != "oidc" {
		return nil
	}

	if OIDCIssuer() == "" {
		return fmt.Errorf("config.ValidateOIDC: OIDC_ISSUER is required")
	}

	if OIDCAuthorizeURL() == "" {
		return fmt.Errorf("config.ValidateOIDC: OIDC_AUTHORIZATION_URL is required")
	}

	if OIDCTokenURL() == "" {
		return fmt.Errorf("config.ValidateOIDC: OIDC_TOKEN_URL is required")
	}

	if OIDCClientID() == "" {
		return fmt.Errorf("config.ValidateOIDC: OIDC_CLIENT_ID is required")
	}

	if OIDCClientSecret() == "" {
		return fmt.Errorf("config.ValidateOIDC: OIDC_CLIENT_SECRET is required")
	}

	if JWTSecret() == "" {
		return fmt.Errorf("config.ValidateOIDC: JWT_SECRET is required")
	}

	if BaseURL() == "" {
		return fmt.Errorf("config.ValidateOIDC: BASE_URL is required")
	}

	return nil
}

func ValidateBasicAuth() error {
	if AuthType() != "basic" {
		return nil
	}

	if BasicAuthUser() == "" {
		return fmt.Errorf("config.ValidateBasicAuth: BASIC_AUTH_USER is required")
	}

	if BasicAuthPassword() == "" {
		return fmt.Errorf("config.ValidateBasicAuth: BASIC_AUTH_PASSWORD is required")
	}

	return nil
}
