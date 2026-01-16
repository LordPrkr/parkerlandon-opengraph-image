package config

import (
	"log/slog"
)

type Server struct {
	CorsAllowedOrigins []string
	ServerCertFile     string
	ServerKeyFile      string
	LogLevel           slog.Leveler
	BaseURL            string
	BrowserPath        string
}

func ParseServerConfig(getEnv func(string) string) *Server {
	profile := getEnv("PROFILE")
	return parseServerConfig(profile)
}
