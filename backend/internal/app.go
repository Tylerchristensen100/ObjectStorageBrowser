package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Tylerchristensen100/object_browser/internal/object_store"
	"github.com/Tylerchristensen100/object_browser/internal/sso"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

const (
	SecretKeyEndpoint  = "ENDPOINT"
	SecretKeyAccessKey = "KEY"
	SecretKeySecretKey = "PASSWORD"
	SecretKeyIssuer    = "ISSUER"
	SecretOauthSecret  = "OAUTH_CLIENT_SECRET"
	SecretOauthKey     = "OAUTH_KEY_ID"
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
		return Config{}, fmt.Errorf("internal/configFromYaml: port must be set in config file")
	}
	if config.Address == "" {
		config.Address = "localhost"
	}

	if config.OAuth.ClientID == "" || config.OAuth.RedirectURL == "" || config.OAuth.AuthURL == "" || config.OAuth.TokenURL == "" {
		return Config{}, fmt.Errorf("internal/configFromYaml: oauth configuration must be set in config file")
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
			Key    string
			Secret string
		}{
			Key:    os.Getenv(SecretOauthKey),
			Secret: os.Getenv(SecretOauthSecret),
		},
	}

	if secrets.ObjectStore.Endpoint == "" || secrets.ObjectStore.AccessKey == "" || secrets.ObjectStore.SecretKey == "" {
		slog.Error("internal/LoadSecrets: Object storage credentials are not set in environment variables")
		os.Exit(1)
	}

	if secrets.SSO.Secret == "" {
		slog.Error("internal/LoadSecrets: SSO Secret is not set in environment variables")
		os.Exit(1)
	}

	slog.Info("Storage endpoint set to " + secrets.ObjectStore.Endpoint)

	return secrets
}

func LoadOauthConfig(app *App) *sso.Auth {
	clientID := app.Config.OAuth.ClientID
	authUrl := app.Config.OAuth.AuthURL
	tokenUrl := app.Config.OAuth.TokenURL
	redirectUrl := app.Config.OAuth.RedirectURL
	scopes := app.Config.OAuth.Scopes

	if clientID == "" || authUrl == "" || tokenUrl == "" || redirectUrl == "" || len(scopes) == 0 {
		app.Logger().Error("internal/LoadOauthConfig: OAuth configuration is incomplete", "Required fields: client_id, auth_url, token_url, redirect_url, scopes",
			slog.String("client_id", clientID),
			slog.String("auth_url", authUrl),
			slog.String("token_url", tokenUrl),
			slog.String("redirect_url", redirectUrl),
			slog.Any("scopes", scopes),
		)
		return nil
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: app.Secrets.SSO.Secret,
		RedirectURL:  redirectUrl,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authUrl,
			TokenURL: tokenUrl,
		},
	}

	auth, err := sso.Init(cfg, app.Config.OAuth.Issuer, app.Secrets.SSO.Key, &app.Log)
	if err != nil {
		app.Logger().Error("internal/LoadOauthConfig: failed to initialize SSO: " + err.Error())
		return nil
	}

	return auth
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

type App struct {
	Log     slog.Logger
	Ctx     context.Context
	Store   *object_store.ObjectStore
	Secrets Secrets
	Config  Config
	Auth    *sso.Auth
}

type Secrets struct {
	ObjectStore struct {
		Endpoint  string
		AccessKey string
		SecretKey string
	}
	SSO struct {
		Key    string
		Secret string
	}
}

type Config struct {
	Address        string      `yaml:"host"`
	Port           int         `yaml:"port"`
	TrustedOrigins []string    `yaml:"trustedOrigins"`
	OAuth          OAuthConfig `yaml:"oauth"`
}

type OAuthConfig struct {
	ClientID    string   `yaml:"client_id"`
	RedirectURL string   `yaml:"redirect_url"`
	Issuer      string   `yaml:"issuer"`
	AuthURL     string   `yaml:"auth_url"`
	TokenURL    string   `yaml:"token_url"`
	Scopes      []string `yaml:"scopes"`
}
