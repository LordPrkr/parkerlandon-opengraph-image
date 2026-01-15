package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/ParkerGits/go-backend-starter/pkg/config"
	"github.com/ParkerGits/go-backend-starter/pkg/httplib"
	"github.com/ParkerGits/go-backend-starter/pkg/ogimage"
)

func Run(ctx context.Context, w io.Writer, getEnv func(string) string) error {
	config := config.ParseServerConfig(getEnv)

	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: config.LogLevel,
	}))
	slog.SetDefault(logger)

	// Initialize browser pool
	browserPool, err := ogimage.NewBrowserPool(config.BrowserPoolSize)
	if err != nil {
		return fmt.Errorf("failed to create browser pool: %w", err)
	}
	defer browserPool.Close()

	// Create generator
	generator := ogimage.NewGenerator(browserPool, config.BaseURL)

	tlsConfig, err := httplib.TLSConfigFromCerts(ctx, config.ServerCertFile, config.ServerKeyFile)
	if err != nil {
		return err
	}

	router := newRouter(config, generator)

	server, err := New(WithPort(3000), WithHandler(router), WithTLSConfig(tlsConfig))
	if err != nil {
		return err
	}

	return server.Run(ctx)
}
