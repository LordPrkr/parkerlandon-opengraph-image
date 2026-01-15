package config

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultProfile = "local"
)

var logLevels = []string{"DEBUG", "INFO", "WARN", "ERROR"}

type serverConfig struct {
	Cors    CorsConfig    `mapStructure:"cors"`
	Log     LogConfig     `mapStructure:"log"`
	Tls     TlsConfig     `mapStructure:"tls"`
	OGImage OGImageConfig `mapstructure:"ogimage"`
}

type CorsConfig struct {
	AllowedOrigins string `mapstructure:"allowed_origins"`
}

type TlsConfig struct {
	Server TlsServerConfig `mapstructure:"server"`
}

type TlsServerConfig struct {
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type OGImageConfig struct {
	BrowserPoolSize int    `mapstructure:"browser_pool_size"`
	BaseURL         string `mapstructure:"base_url"`
}

func parseServerConfig(profile string) *Server {
	if profile == "" {
		profile = defaultProfile
	}
	viper.SetConfigName(profile)
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs/server/")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config: %w", err))
	}

	var config serverConfig
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	fmt.Printf("%#v\n", config)

	return mapConfig(config)
}

func mapConfig(config serverConfig) *Server {
	return &Server{
		CorsAllowedOrigins: parseCorsAllowedOrigins(config.Cors.AllowedOrigins),
		ServerCertFile:     config.Tls.Server.CertFile,
		ServerKeyFile:      config.Tls.Server.KeyFile,
		LogLevel:           parseLogLevel(config.Log.Level),
		BrowserPoolSize:    config.OGImage.BrowserPoolSize,
		BaseURL:            config.OGImage.BaseURL,
	}
}

func parseLogLevel(level string) slog.Leveler {
	if !slices.Contains(logLevels, level) {
		panic(fmt.Sprintf("expected config value log.level to be one of %v, received %s", logLevels, level))
	}
	slogLevel := &slog.LevelVar{}
	slogLevel.UnmarshalText([]byte(level))
	return slogLevel
}

func parseCorsAllowedOrigins(allowedOrigins string) []string {
	return strings.Split(allowedOrigins, ",")
}
