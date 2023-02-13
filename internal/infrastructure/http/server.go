package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"

	"github.com/aplulu/gcsproxy/internal/config"
)

var server http.Server

func RunServer() error {
	serverCtx := context.Background()
	storageClient, err := storage.NewClient(serverCtx)
	if err != nil {
		return fmt.Errorf("http.RunServer: failed to create storage client: %w", err)
	}
	storageBucket := storageClient.Bucket(config.GoogleCloudStorageBucket())

	httpMux := chi.NewRouter()

	httpMux.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		key := req.URL.Path
		if len(config.MainPageSuffix()) > 0 && strings.HasSuffix(key, "/") {
			key += config.MainPageSuffix()
		}

		key = strings.TrimPrefix(key, "/")

		obj := storageBucket.Object(key)
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			responseError(w, err)
			return
		}

		r, err := obj.NewReader(ctx)
		if err != nil {
			responseError(w, err)
			return
		}

		writeHeaders(w, attrs, r.Attrs)

		if _, err := io.Copy(w, r); err != nil {
			log.Printf("http.RunServer: failed to copy content: %v\n", err)
		}
	})

	server = http.Server{
		Addr:    net.JoinHostPort(config.Listen(), config.Port()),
		Handler: httpMux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func writeHeaders(w http.ResponseWriter, attrs *storage.ObjectAttrs, attrs2 storage.ReaderObjectAttrs) {
	writeStringHeader(w, "Last-Modified", attrs.Updated.Format(http.TimeFormat))
	writeStringHeader(w, "Content-Type", attrs.ContentType)
	writeInt64Header(w, "Content-Length", attrs.Size)
	writeStringHeader(w, "Content-Encoding", attrs.ContentEncoding)
	writeStringHeader(w, "Content-Disposition", attrs.ContentDisposition)
	writeStringHeader(w, "Cache-Control", attrs.CacheControl)
}

func responseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrObjectNotExist):
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func StopServer(ctx context.Context) error {
	return server.Shutdown(ctx)
}
