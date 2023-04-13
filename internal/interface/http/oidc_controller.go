package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"

	"github.com/aplulu/gcsproxy/internal/config"
	"github.com/aplulu/gcsproxy/internal/domain/model"
	"github.com/aplulu/gcsproxy/internal/util"
)

const (
	oidcSessionCookieName = "_gpso"
	oidcSessionTTL        = 300
	authSessionCookieName = "_gpsa"
)

type OIDCController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}

type oidcController struct {
}

// Login is the handler for the login route.
func (c *oidcController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get OIDC config
	oc, err := model.GetOIDCConfig(ctx)
	if err != nil {
		responseError(w, err)
		return
	}

	// Create auth session
	sessStr, sess, err := model.NewOIDCSession(r.URL.Query().Get("redirect"))
	if err != nil {
		responseError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oidcSessionCookieName,
		Value:    sessStr,
		Path:     "/",
		Secure:   util.IsTLS(r),
		MaxAge:   oidcSessionTTL,
		HttpOnly: true,
	})

	var authorizeURL string
	if config.OIDCProvider() == "https://accounts.google.com" && config.OIDCGoogleHostedDomain() != "" {
		authorizeURL = oc.AuthCodeURL(sess.State, oauth2.SetAuthURLParam("hd", config.OIDCGoogleHostedDomain()))
	} else {
		authorizeURL = oc.AuthCodeURL(sess.State)
	}

	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

// Callback is the handler for the callback route.
func (c *oidcController) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get OIDC session
	sessCookie, err := r.Cookie(oidcSessionCookieName)
	if err != nil {
		responseError(w, err)
		return
	}
	sess, err := model.ParseOIDCSession(sessCookie.Value)
	if err != nil {
		responseError(w, err)
		return
	}
	if sess.State != r.URL.Query().Get("state") {
		responseError(w, model.ErrInvalidState)
		return
	}

	// Exchange code for token
	token, err := model.ExchangeOIDCToken(ctx, r.URL.Query().Get("code"))
	if err != nil {
		responseError(w, err)
		return
	}

	// Create Auth session
	sessToken, exp, err := model.CreateAuthSession(token.Sub)
	if err != nil {
		responseError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authSessionCookieName,
		Value:    sessToken,
		Path:     "/",
		Secure:   util.IsTLS(r),
		Expires:  *exp,
		HttpOnly: true,
	})
	http.Redirect(w, r, config.BaseURL()+sess.RedirectURL, http.StatusFound)
}

func NewOIDCController() OIDCController {
	return &oidcController{}
}

func Register(mux *chi.Mux) {
	controller := NewOIDCController()

	mux.Get("/login", controller.Login)
	mux.Get("/callback", controller.Callback)
}
