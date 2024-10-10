package common

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// AppConfig is the structured configuration used throughout the application.
type AppConfig struct {
	App  AppSettings `mapstructure:"app"`
	DB   DBConfig    `mapstructure:"db"`
	JWT  JWTConfig   `mapstructure:"jwt"`
	Card CardConfig  `mapstructure:"card"`
}

type AppSettings struct {
	Env           string `mapstructure:"env"`
	GinMode       string `mapstructure:"gin_mode"`
	ServerAddress string `mapstructure:"server_address"`
}

type DBConfig struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

type JWTConfig struct {
	PrivateKey        string `mapstructure:"private_key"`
	PublicKey         string `mapstructure:"public_key"`
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

type CardConfig struct {
	AESKey string `mapstructure:"aes_key"`
}

// LoadConfig reads the config file and returns a structured AppConfig.
func LoadConfig() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	bindEnvVariables(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set default values for JWT expiration times
	config.JWT.AccessExpiration = 30 * time.Minute
	config.JWT.RefreshExpiration = 24 * time.Hour

	if err := decodeKeys(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// decodeKeys decodes all base64 keys
func decodeKeys(config *AppConfig) error {
	decodeBase64 := func(encoded string) (string, error) {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 string: %w", err)
		}
		return string(decoded), nil
	}

	var err error
	config.JWT.PrivateKey, err = decodeBase64(config.JWT.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode JWT private key: %w", err)
	}

	config.JWT.PublicKey, err = decodeBase64(config.JWT.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode JWT public key: %w", err)
	}

	config.Card.AESKey, err = decodeBase64(config.Card.AESKey)
	if err != nil {
		return fmt.Errorf("failed to decode Card AES key: %w", err)
	}

	return nil
}

// validateConfig checks if all required fields are set in the AppConfig.
func validateConfig(config *AppConfig) error {
	checks := []struct {
		name  string
		valid bool
	}{
		{"app.env", config.App.Env != ""},
		{"app.server_address", config.App.ServerAddress != ""},
		{"app.gin_mode", config.App.GinMode != ""},
		{"db.url", config.DB.URL != ""},
		{"db.max_open_conns", config.DB.MaxOpenConns > 0},
		{"db.max_idle_conns", config.DB.MaxIdleConns > 0},
		{"db.conn_max_lifetime", config.DB.ConnMaxLifetime > 0},
		{"db.conn_max_idle_time", config.DB.ConnMaxIdleTime > 0},
		{"jwt.private_key", config.JWT.PrivateKey != ""},
		{"jwt.public_key", config.JWT.PublicKey != ""},
		{"card.aes_key", config.Card.AESKey != ""},
	}

	var missingConfigs []string
	for _, check := range checks {
		if !check.valid {
			missingConfigs = append(missingConfigs, check.name)
		}
	}

	if len(missingConfigs) > 0 {
		return fmt.Errorf("missing or invalid required configurations: %v", missingConfigs)
	}

	return nil
}

// bindEnvVariables maps environment variables to configuration keys.
// This allows overriding config values using Docker environment variables.
//
// Docker example:
//
//	docker run -e DB_URL="postgres://user:pass@host:5432/db" -e APP_ENV="production" ...
func bindEnvVariables(v *viper.Viper) {
	envMappings := map[string]string{
		"app.env":               "APP_ENV",
		"app.gin_mode":          "GIN_MODE",
		"app.server_address":    "SERVER_ADDRESS",
		"db.url":                "DB_URL",
		"db.max_open_conns":     "DB_MAX_OPEN_CONNS",
		"db.max_idle_conns":     "DB_MAX_IDLE_CONNS",
		"db.conn_max_lifetime":  "DB_CONN_MAX_LIFETIME",
		"db.conn_max_idle_time": "DB_CONN_MAX_IDLE_TIME",
		"jwt.private_key":       "JWT_PRIVATE_KEY",
		"jwt.public_key":        "JWT_PUBLIC_KEY",
		"card.aes_key":          "CARD_AES_KEY",
	}

	for configKey, envVar := range envMappings {
		v.BindEnv(configKey, envVar)
	}
}
