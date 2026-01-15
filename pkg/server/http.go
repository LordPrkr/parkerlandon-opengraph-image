package server

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

const (
	defaultHTTPPort = "3000"
)

type HttpServer struct {
	httpServer *http.Server
}

type options struct {
	port      string
	handler   http.Handler
	tlsConfig *tls.Config
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(o *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		o.port = ":" + strconv.Itoa(port)
		return nil
	}
}

func WithHandler(handler http.Handler) Option {
	return func(o *options) error {
		o.handler = handler
		return nil
	}
}

func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(o *options) error {
		o.tlsConfig = tlsConfig
		return nil
	}
}

func New(opts ...Option) (*HttpServer, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	addr := defaultHTTPPort
	if options.port != "" {
		addr = options.port
	}

	httpServer := http.Server{
		Addr:      addr,
		Handler:   options.handler,
		ErrorLog:  slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
		TLSConfig: options.tlsConfig,
	}
	server := &HttpServer{
		httpServer: &httpServer,
	}

	return server, nil

}

func (s *HttpServer) Run(ctx context.Context) error {
	logger := slog.Default()
	if s.httpServer.TLSConfig != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"starting server with tls",
			slog.String("address", s.httpServer.Addr),
		)
		return s.httpServer.ListenAndServeTLS("", "")
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "starting server", slog.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}
