package sso

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const (
	authCookieKey = "auth_token"
)

type Auth struct {
	Config *oauth2.Config
	oidc   *oidcConfig
	ctx    context.Context
	log    *slog.Logger
	client http.Client
}

type oidcConfig struct {
	OIDC     *oidc.Provider
	Verifier oidc.IDTokenVerifier
	key      *rsa.PrivateKey
	keyId    string
}

func Init(config *oauth2.Config, issuer string, keyID string, log *slog.Logger) (*Auth, error) {
	var usingOidc bool = true
	if issuer == "" {
		log.Warn("No OIDC issuer provided, falling back to generic OAuth2")
		usingOidc = false
	}

	if usingOidc {
		provider, err := oidc.NewProvider(context.Background(), issuer)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				log.Error("sso/Init: OIDC issuer not found (404). Please check the issuer URL.", "issuer", issuer)
				return nil, errors.New(errorsOidcIssuerNotFound)
			}
			return nil, err
		}
		config.Endpoint = provider.Endpoint()
		config.Scopes = append([]string{oidc.ScopeOpenID}, config.Scopes...)

		verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

		privateKey, err := loadPrivateKey(config.ClientSecret)
		if err != nil {
			log.Error("sso/Init: Failed to load private key", "error", err)
			return nil, err
		}
		client := http.DefaultClient
		auth := &Auth{
			Config: config,
			ctx:    context.Background(),
			log:    log,
			oidc:   &oidcConfig{OIDC: provider, Verifier: *verifier, key: privateKey, keyId: keyID},
			client: *client,
		}
		return auth, nil
	} else {
		auth := &Auth{
			Config: config,
			ctx:    context.Background(),
			log:    log,
			oidc:   nil,
		}
		return auth, nil
	}
}

func (a *Auth) SetAuthCookie(token *oauth2.Token, w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieKey,
		Value:    token.AccessToken,
		HttpOnly: true,
		Secure:   true,
	})
	return nil

}

func loadPrivateKey(keyData string) (*rsa.PrivateKey, error) {
	var key *rsa.PrivateKey
	keyBytes := []byte(keyData)

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing RSA private key")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		var err error

		// PKCS1
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	case "PRIVATE KEY":
		var ok bool
		// PKCS8
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		key, ok = k.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}

	}
	return key, nil
}
