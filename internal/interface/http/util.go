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
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
