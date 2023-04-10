package middleware

import (
	"fmt"
	"net/http"

	"github.com/aplulu/gcsproxy/pkg/accesstoken"
)

// AuthOIDCConfig is the configuration for the Auth middleware.
type AuthOIDCConfig struct {
	CookieName  string
	Issuer      string
	Audience    string
	SecretKey   string
	RedirectURL string
	Skipper     Skipper
}

// AuthOIDCWithConfig returns a middleware that authenticates requests.
func AuthOIDCWithConfig(conf AuthOIDCConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if conf.Skipper != nil && conf.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			gps, _ := r.Cookie(conf.CookieName)
			if gps != nil {
				_, err := accesstoken.ParseAccessToken(gps.Value, conf.Issuer, conf.Audience, []byte(conf.SecretKey))
				if err == nil {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("Location", fmt.Sprintf("%s?redirect=%s", conf.RedirectURL, r.URL.Path))
			w.WriteHeader(http.StatusFound)
		})
	}
}
