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

	// Initialize browser manager
	browserManager, err := ogimage.NewBrowserManager(config.BrowserPath)
	if err != nil {
		return fmt.Errorf("failed to create browser manager: %w", err)
	}
	defer browserManager.Close()

	// Create generator
	generator := ogimage.NewGenerator(browserManager, config.BaseURL)

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
