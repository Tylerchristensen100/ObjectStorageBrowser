package sso

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (a *Auth) Token(code string) (*oauth2.Token, error) {
	unique := uuid.New().String()
	claims := jwt.MapClaims{
		"iss": "343507177119483263",
		"sub": "343507177119483263",
		"aud": a.Config.Endpoint.TokenURL,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iat": time.Now().Unix(),
		"jti": unique,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = a.oidc.keyId

	signedJWT, err := token.SignedString(a.oidc.key) // privateKey
	if err != nil {
		return nil, err
	}

	// Step 2: Build token request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", a.Config.RedirectURL)
	data.Set("client_id", a.Config.ClientID)
	data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	data.Set("client_assertion", signedJWT)
	a.log.Info("Signed JWT", "jwt", signedJWT)

	req, err := http.NewRequest(http.MethodPost, a.Config.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req.WithContext(a.ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		a.log.Error("Token exchange failed", "status", resp.Status, "body", bodyString)

		return nil, errors.New(errorTokenExchangeFailed + " status: " + resp.Status)
	}

	var tokenResp oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (a *Auth) validateToken(token string) (*oidc.IDToken, error) {
	id, err := a.oidc.Verifier.Verify(a.ctx, token)
	if err != nil {
		return nil, err
	}
	return id, nil
}
func (a *Auth) User(id *oidc.IDToken) (*Claims, error) {
	var claims Claims
	if err := id.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (a *Auth) Validate(r *http.Request) (bool, *oidc.IDToken, error) {
	cookie, err := r.Cookie(authCookieKey)
	if err != nil {
		return false, nil, err
	}

	id, err := a.validateToken(cookie.Value)
	if err != nil {
		return false, nil, err
	}

	if id == nil {
		return false, nil, errors.New(errorTokenInvalid)
	}

	if id.VerifyAccessToken(cookie.Value) != nil {
		return false, nil, nil
	}

	a.log.Debug("Validated token for subject", "sub", id.Subject)
	return true, id, nil
}
