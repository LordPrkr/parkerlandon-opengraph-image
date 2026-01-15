package httplib

import (
	"context"
	"crypto/tls"
	"log/slog"
)

func TLSConfigFromCerts(ctx context.Context, certFile, keyFile string) (*tls.Config, error) {
	logger := slog.Default()

	if certFile == "" || keyFile == "" {
		return nil, nil
	}

	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"reading tls certs",
		slog.String("cert", certFile),
		slog.String("key", keyFile),
	)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}
