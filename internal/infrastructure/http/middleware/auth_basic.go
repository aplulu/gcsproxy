package middleware

import (
	"crypto/subtle"
	"net/http"
)

// AuthBasicConfig is the configuration for the AuthBasic middleware.
type AuthBasicConfig struct {
	User     string
	Password string
	Skipper  Skipper
}

// AuthBasicWithConfig returns a middleware that authenticates requests.
func AuthBasicWithConfig(conf AuthBasicConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if conf.Skipper != nil && conf.Skipper(r) {
				next.ServeHTTP(w, r)
				return
			}

			user, pass, ok := r.BasicAuth()
			if ok && subtle.ConstantTimeCompare([]byte(user), []byte(conf.User)) == 1 && subtle.ConstantTimeCompare([]byte(pass), []byte(conf.Password)) == 1 {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", "Basic realm=\"GCS Proxy\"")
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
