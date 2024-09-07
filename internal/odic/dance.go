package oidc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrAuthenticationRequestFailed = errors.New("authentication request failed")
	ErrMissingCode                 = errors.New("missing code")
	ErrMissingState                = errors.New("missing state")
	ErrTokenRequestFailed          = errors.New("token request failed")
	ErrInvalidIdToken              = errors.New("id_token validation failed")
	ErrInvalidSession              = errors.New("invalid session")
)

func GenerateAuthenticationRequestUrl(clientId string, redirectUri string, kvs kvs) (redirectUrl, state string, err error) {
	d := getDiscovery()
	u, err := url.Parse(d.AuthorizationEndpoint)
	if err != nil {
		return "", "", err
	}

	state, err = RandomString(16)
	if err != nil {
		return "", "", err
	}
	nonce, err := RandomString(16)
	if err != nil {
		return "", "", err
	}

	q := u.Query()
	q.Add("scope", "openid")
	q.Add("response_type", "code")
	q.Add("client_id", clientId)
	q.Add("redirect_uri", redirectUri)
	q.Add("state", state)
	q.Add("nonce", nonce)

	u.RawQuery = q.Encode()

	kvs.SetNonce(state, nonce, 600)

	return u.String(), state, nil
}

func CallbackCode(
	state string,
	callbackQuery map[string][]string,
	clientId, clientSecret, redirectUri string,
	kvs kvs,
	aud string,
) (*JWTPayload, error) {
	toQueryMap := func(q map[string][]string) map[string]string {
		m := make(map[string]string)

		for k, v := range q {
			m[k] = v[0]
		}
		return m
	}

	q := toQueryMap(callbackQuery)

	qerr, ok := q["error"]
	if ok {
		return nil, errors.New(qerr)
	}

	code, ok := q["code"]
	if !ok {
		return nil, ErrMissingCode
	}
	returnedState, ok := q["state"]
	if !ok {
		return nil, ErrMissingState
	}

	// state の検証
	if state != returnedState {
		return nil, fmt.Errorf("state が違う expect:%s, got:%s", state, returnedState)
	}

	nonce, err := kvs.GetNonce(state)
	if err != nil {
		return nil, fmt.Errorf("nonce がない %w", err)
	}
	defer func() {
		kvs.DeleteNonce(state)
	}()

	tokenResponse, err := tokenRequest(code, clientId, clientSecret, redirectUri)
	if err != nil {
		return nil, ErrTokenRequestFailed
	}

	idToken := tokenResponse.IdToken
	payload, err := validateIdToken(idToken, nonce, aud)
	if err != nil {
		return nil, ErrInvalidIdToken
	}

	return payload, nil
}

func validateIdToken(idToken, nonce, aud string) (*JWTPayload, error) {
	jwk := getJWK()
	if len(jwk.Keys) != 1 {
		return nil, ErrSingleKeyIsSupported
	}

	if err := verifyIDToken(idToken, jwk.Keys[0]); err != nil {
		return nil, err
	}

	dis := getDiscovery()

	claims, err := decodeIDToken(idToken)
	if err != nil {
		return nil, err
	}

	payload := claims.Payload
	if payload.Iss != dis.Issuer {
		return nil, errors.New("invalid issuer")
	}
	if payload.Aud != aud {
		return nil, errors.New("invalid aud")
	}
	if payload.Nonce != nonce {
		return nil, errors.New("invalid nonce")
	}

	now := uint64(time.Now().Unix())
	if now > payload.Exp {
		return nil, errors.New("expired token")
	}
	if now < payload.Iat {
		return nil, errors.New("invalid iat")
	}

	return &payload, nil
}

type TokenResponse struct {
	IdToken string `json:"id_token"`
}

func tokenRequest(code, clientId, clientSecret, redirectUri string) (*TokenResponse, error) {
	q := make(url.Values)

	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("client_id", clientId)
	q.Add("client_secret", clientSecret)
	q.Add("redirect_uri", redirectUri)

	d := getDiscovery()
	req, err := http.NewRequest(http.MethodPost, d.TokenEndpoint, strings.NewReader(q.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(resb, &tokenResponse); err != nil {
		return nil, err
	}
	// 明らかになんか短い
	if len(tokenResponse.IdToken) <= 5 {
		return nil, errors.New("invalid TokenEndpoint response format")
	}

	return &tokenResponse, nil
}
