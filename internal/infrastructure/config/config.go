package config

import (
	"fmt"
	"os"
	"time"
)

type Config interface {
	GetAppConfig() AppConfig
	GetDatabaseConfig() DatabaseConfig
	GetLoggingConfig() LoggingConfig
	GetServerConfig() ServerConfig
	GetCORSConfig() CORSConfig
}

type AppConfig struct {
	Name        string `mapstructure:"name" validate:"required" default:"Todo API"`
	Version     string `mapstructure:"version" validate:"required" default:"1.0.0"`
	Environment string `mapstructure:"environment" validate:"required,oneof=development test staging production local" default:"development"`
	Port        int    `mapstructure:"port" validate:"required,min=1,max=65535" default:"8080"`
	Host        string `mapstructure:"host" validate:"required" default:"localhost"`
}

func (a AppConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

type DatabaseConfig struct {
	Driver            string        `mapstructure:"driver" validate:"required,oneof=postgres sqlite" default:"postgres"`
	Host              string        `mapstructure:"host" validate:"required" default:"localhost"`
	Port              int           `mapstructure:"port" validate:"required,min=1,max=65535" default:"5432"`
	Name              string        `mapstructure:"name" validate:"required" default:"todo_db"`
	Username          string        `mapstructure:"user" validate:"required" default:"postgres"`
	Password          string        `mapstructure:"password" validate:"required" default:"postgres"`
	ConnectionURL     string        `mapstructure:"url" validate:"omitempty"`
	SSLMode           string        `mapstructure:"ssl_mode" validate:"required,oneof=disable require verify-ca verify-full" default:"disable"`
	MaxOpenConns      int           `mapstructure:"max_open_conns" validate:"required,min=1,max=100" default:"25"`
	MaxIdleConns      int           `mapstructure:"max_idle_conns" validate:"required,min=1,max=50" default:"5"`
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime" validate:"required" default:"5m"`
}

func (d DatabaseConfig) GetConnectionURL() string {
	if d.ConnectionURL != "" {
		return d.ConnectionURL
	}
	
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		d.Driver,
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

type LoggingConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error fatal" default:"info"`
	Format string `mapstructure:"format" validate:"required,oneof=json console" default:"json"`
	Output string `mapstructure:"output" validate:"required" default:"stdout"`
}

type ServerConfig struct {
	ReadTimeout           time.Duration `mapstructure:"read_timeout" validate:"required" default:"30s"`
	WriteTimeout          time.Duration `mapstructure:"write_timeout" validate:"required" default:"30s"`
	IdleTimeout           time.Duration `mapstructure:"idle_timeout" validate:"required" default:"60s"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout" validate:"required" default:"10s"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins" validate:"required,min=1" default:"[\"*\"]"`
	AllowedMethods   []string `mapstructure:"allowed_methods" validate:"required,min=1" default:"[\"GET\",\"POST\",\"PUT\",\"DELETE\",\"OPTIONS\"]"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" validate:"required,min=1" default:"[\"*\"]"`
	AllowCredentials bool     `mapstructure:"allow_credentials" default:"true"`
}

type viperConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Server   ServerConfig   `mapstructure:"server"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

func NewConfig() (Config, error) {
	configFileName := os.Getenv("CONFIG_FILE_NAME")
	if configFileName == "" {
		configFileName = "config.local"
	}
	return loadConfig(configFileName)
}

func (c *viperConfig) GetAppConfig() AppConfig {
	return c.App
}

func (c *viperConfig) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}

func (c *viperConfig) GetLoggingConfig() LoggingConfig {
	return c.Logging
}

func (c *viperConfig) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *viperConfig) GetCORSConfig() CORSConfig {
	return c.CORS
} 