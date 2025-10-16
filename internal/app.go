package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Tylerchristensen100/object_browser/internal/object_store"
	"gopkg.in/yaml.v3"
)

type App struct {
	Log     slog.Logger
	Ctx     context.Context
	Store   *object_store.ObjectStore
	Secrets Secrets
	Config  Config
}

type Secrets struct {
	ObjectStore struct {
		Endpoint  string
		AccessKey string
		SecretKey string
	}
	SSO struct {
		Issuer   string
		ClientID string
	}
}

type Config struct {
	Address        string   `yaml:"host"`
	Port           int      `yaml:"port"`
	TrustedOrigins []string `yaml:"trustedOrigins"`
}

const (
	SecretKeyEndpoint  = "ENDPOINT"
	SecretKeyAccessKey = "KEY"
	SecretKeySecretKey = "PASSWORD"
	SecretKeyIssuer    = "ISSUER"
	SecretKeyClientID  = "CLIENT_ID"
)

func ConfigFromYaml(logger *slog.Logger, path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		logger.Error("internal/configFromYaml: " + err.Error())
		return Config{}, err
	}
	defer file.Close()

	var config Config

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Error("internal/configFromYaml: " + err.Error())
		return Config{}, err
	}

	if config.Port == 0 {
		return Config{}, fmt.Errorf("port must be set in config file")
	}
	if config.Address == "" {
		config.Address = "localhost"
	}

	return config, nil
}

func LoadSecrets() Secrets {
	secrets := Secrets{
		ObjectStore: struct {
			Endpoint  string
			AccessKey string
			SecretKey string
		}{
			Endpoint:  os.Getenv(SecretKeyEndpoint),
			AccessKey: os.Getenv(SecretKeyAccessKey),
			SecretKey: os.Getenv(SecretKeySecretKey),
		},
		SSO: struct {
			Issuer   string
			ClientID string
		}{
			Issuer:   os.Getenv(SecretKeyIssuer),
			ClientID: os.Getenv(SecretKeyClientID),
		},
	}

	if secrets.ObjectStore.Endpoint == "" || secrets.ObjectStore.AccessKey == "" || secrets.ObjectStore.SecretKey == "" {
		slog.Error("Object storage credentials are not set in environment variables")
		os.Exit(1)
	}

	if secrets.SSO.Issuer == "" || secrets.SSO.ClientID == "" {
		slog.Error("SSO credentials are not set in environment variables")
		os.Exit(1)
	}

	return secrets
}

func (c *Config) Location() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

func (a *App) Logger() *slog.Logger {
	return &a.Log
}

func (a *App) Context() context.Context {
	return a.Ctx
}
