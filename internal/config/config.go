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
	OIDCProvider             string   `envconfig:"oidc_provider" default:"https://accounts.google.com"`
	OIDCScopes               []string `envconfig:"oidc_scopes" default:"openid,profile,email"`
	OIDCAuthorizeURL         string   `envconfig:"oidc_authorize_url" default:""`
	OIDCTokenURL             string   `envconfig:"oidc_token_url" default:""`
	OIDCClientID             string   `envconfig:"oidc_client_id" default:""`
	OIDCClientSecret         string   `envconfig:"oidc_client_secret" default:""`
	OIDCGoogleHostedDomain   string   `envconfig:"oidc_google_hosted_domain" default:""`
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

func OIDCProvider() string {
	return conf.OIDCProvider
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

func OIDCGoogleHostedDomain() string {
	return conf.OIDCGoogleHostedDomain
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

	if OIDCProvider() == "" {
		return fmt.Errorf("config.ValidateOIDC: OIDC_PROVIDER is required")
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
