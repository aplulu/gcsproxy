package http

import (
	"errors"
	"net/http"

	"github.com/aplulu/gcsproxy/internal/domain/model"
)

func responseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidRedirectURL):
	case errors.Is(err, model.ErrInvalidState):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, model.ErrInvalidHostedDomain):
		http.Error(w, "Access with this Google account is not allowed", http.StatusForbidden)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
