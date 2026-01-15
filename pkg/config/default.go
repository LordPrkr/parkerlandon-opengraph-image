package config

import (
	"log/slog"
)

type Server struct {
	CorsAllowedOrigins []string
	ServerCertFile     string
	ServerKeyFile      string
	LogLevel           slog.Leveler
	BrowserPoolSize    int
	BaseURL            string
}

func ParseServerConfig(getEnv func(string) string) *Server {
	profile := getEnv("PROFILE")
	return parseServerConfig(profile)
}
