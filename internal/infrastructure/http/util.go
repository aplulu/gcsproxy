package http

import (
	"errors"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"

	"github.com/aplulu/gcsproxy/internal/config"
	"github.com/aplulu/gcsproxy/internal/domain/model"
)

func writeHeaders(w http.ResponseWriter, attrs *storage.ObjectAttrs, chunked bool) {
	writeStringHeader(w, "Last-Modified", attrs.Updated.Format(http.TimeFormat))
	writeStringHeader(w, "Content-Type", attrs.ContentType)
	writeStringHeader(w, "Content-Disposition", attrs.ContentDisposition)
	writeStringHeader(w, "Content-Encoding", attrs.ContentEncoding)

	if chunked {
		writeStringHeader(w, "Transfer-Encoding", "chunked")
	} else {
		writeInt64Header(w, "Content-Length", attrs.Size)
	}

	// do not cache if authentication is enabled
	if config.AuthType() == "oidc" || config.AuthType() == "basic" {
		writeStringHeader(w, "Cache-Control", "private, max-age=60")
	} else {
		writeStringHeader(w, "Cache-Control", attrs.CacheControl)
	}
}

func responseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrObjectNotExist):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, model.ErrStreamingUnsupported):
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeStringHeader(w http.ResponseWriter, name string, value string) {
	if len(value) > 0 {
		w.Header().Set(name, value)
	}
}

func writeInt64Header(w http.ResponseWriter, name string, value int64) {
	if value > 0 {
		w.Header().Set(name, strconv.FormatInt(value, 10))
	}
}
