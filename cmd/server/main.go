package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ParkerGits/go-backend-starter/pkg/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Stdout, os.Getenv); err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}
}
