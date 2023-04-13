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
	"github.com/aplulu/gcsproxy/internal/domain/model"
	"github.com/aplulu/gcsproxy/internal/infrastructure/http/middleware"
	appHttp "github.com/aplulu/gcsproxy/internal/interface/http"
)

const (
	gcsProxyPathPrefix = "/_gcsproxy"
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

	// OpenID Connect
	if config.AuthType() == "oidc" {
		if err := config.ValidateOIDC(); err != nil {
			return fmt.Errorf("http.RunServer: invalid OIDC config: %w", err)
		}

		httpMux.Use(middleware.AuthOIDCWithConfig(middleware.AuthOIDCConfig{
			CookieName:  "_gpsa",
			Issuer:      config.BaseURL(),
			Audience:    config.BaseURL(),
			SecretKey:   config.JWTSecret(),
			RedirectURL: config.BaseURL() + gcsProxyPathPrefix + "/oidc/login",
			Skipper: func(r *http.Request) bool {
				return strings.HasPrefix(r.URL.Path, gcsProxyPathPrefix)
			},
		}))

		authMux := chi.NewRouter()
		appHttp.Register(authMux)
		httpMux.Mount(gcsProxyPathPrefix+"/oidc", authMux)
	} else if config.AuthType() == "basic" { // Basic Auth
		if err := config.ValidateBasicAuth(); err != nil {
			return fmt.Errorf("http.StartServer: invalid Basic Auth config: %w", err)
		}

		httpMux.Use(middleware.AuthBasicWithConfig(middleware.AuthBasicConfig{
			User:     config.BasicAuthUser(),
			Password: config.BasicAuthPassword(),
			Skipper: func(r *http.Request) bool {
				return strings.HasPrefix(r.URL.Path, gcsProxyPathPrefix)
			},
		}))
	}

	httpMux.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, gcsProxyPathPrefix) {
			http.NotFound(w, req)
			return
		}

		ctx := req.Context()

		key := req.URL.Path
		if len(config.MainPageSuffix()) > 0 && strings.HasSuffix(key, "/") {
			key += config.MainPageSuffix()
		}

		key = strings.TrimPrefix(key, "/")

		serveFile(ctx, storageBucket, key, w)
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

func serveFile(ctx context.Context, storageBucket *storage.BucketHandle, key string, w http.ResponseWriter) {
	obj := storageBucket.Object(key)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		// fallback to Not Found Page
		if errors.Is(err, storage.ErrObjectNotExist) && len(config.NotFoundPage()) > 0 && key != config.NotFoundPage() {
			serveFile(ctx, storageBucket, config.NotFoundPage(), w)
			return
		}
		responseError(w, err)
		return
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		responseError(w, err)
		return
	}

	// if file size is larger than 32MB, use chunked transfer encoding
	if attrs.Size > 32*1024*1024 {
		flusher, ok := w.(http.Flusher)
		if !ok {
			responseError(w, model.ErrStreamingUnsupported)
			return
		}

		// write headers
		writeHeaders(w, attrs, true)

		buf := make([]byte, 32*1024)
		for {
			n, err := r.Read(buf)
			if err != nil {
				if err == io.EOF || err == context.Canceled {
					break
				}
				log.Printf("http.RunServer: failed to read content: %v\n", err)
				break
			}
			if _, err := w.Write(buf[:n]); err != nil {
				log.Printf("http.RunServer: failed to write content: %v\n", err)
				break
			}
			flusher.Flush()
		}
	} else {
		// write headers
		writeHeaders(w, attrs, false)

		if _, err := io.Copy(w, r); err != nil {
			log.Printf("http.RunServer: failed to copy content: %v\n", err)
		}
	}
}

func StopServer(ctx context.Context) error {
	return server.Shutdown(ctx)
}
