package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`

	Log struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"log"`

	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Jwt struct {
		AccessSecret  string
		RefreshSecret string
	}
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("db.sslmode", "disable")

	viper.SetEnvPrefix("UPTIME")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed read config %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed decode config %w", err)
	}

	cfg.DB.Password = viper.GetString("db.password")
	if cfg.DB.Password == "" {
		return nil, errors.New("password not set in UPTIME_DB_PASSWORD")
	}

	cfg.Jwt.AccessSecret = viper.GetString("jwt.accesssecret")
	if cfg.Jwt.AccessSecret == "" {
		return nil, errors.New("password not set in UPTIME_JWT_ACCESSSECRET")
	}

	cfg.Jwt.RefreshSecret = viper.GetString("jwt.refreshsecret")
	if cfg.Jwt.RefreshSecret == "" {
		return nil, errors.New("password not set in UPTIME_JWT_REFRESHSECRET")
	}

	return &cfg, nil
}
