package oidc

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type JWTHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type JWTPayload struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Exp   uint64 `json:"exp"`
	Iat   uint64 `json:"iat"`
	Nonce string `json:"nonce"`

	// Custom claim
	Roles []string `json:"roles"`
}

type JWT struct {
	Header  JWTHeader
	Payload JWTPayload
}

var (
	ErrInvalidJWTFormat    = errors.New("invalid jwt format")
	ErrUnsupportedAlg      = errors.New("unsupported alg")
	ErrDecodeJWKSigningKey = errors.New("failed tot decode jwk signing key")
	ErrInvalidSignature    = errors.New("invalid signature")
)

func decodeIDToken(idToken string) (*JWT, error) {
	sp := strings.Split(idToken, ".")
	if len(sp) != 3 {
		return nil, ErrInvalidJWTFormat
	}

	hb, err := base64.RawStdEncoding.DecodeString(sp[0])
	if err != nil {
		return nil, ErrInvalidJWTFormat
	}

	pb, err := base64.RawStdEncoding.DecodeString(sp[1])
	if err != nil {
		return nil, ErrInvalidJWTFormat
	}

	var header JWTHeader
	if err := json.Unmarshal(hb, &header); err != nil {
		log.Println(err)
		return nil, ErrInvalidJWTFormat
	}

	var payload JWTPayload
	if err := json.Unmarshal(pb, &payload); err != nil {
		log.Println(err)
		return nil, ErrInvalidJWTFormat
	}

	return &JWT{
		Header:  header,
		Payload: payload,
	}, nil
}

func verifyIDToken(idToken string, key JwkKey) error {
	dec, err := decodeIDToken(idToken)
	if err != nil {
		return ErrInvalidJWTFormat
	}

	if dec.Header.Alg != "RS256" {
		return ErrUnsupportedAlg
	}

	pk, err := key.RsaPubkey()
	if err != nil {
		return ErrDecodeJWKSigningKey
	}

	sp := strings.Split(idToken, ".")

	sig, err := base64.RawURLEncoding.DecodeString(sp[2])
	if err != nil {
		return ErrInvalidSignature
	}

	if err := verifyRS256(sp[0]+"."+sp[1], sig, pk); err != nil {
		return err
	}

	return nil
}

func verifyRS256(msg string, sig []byte, pubkey *rsa.PublicKey) error {
	hasher := sha256.New()
	hasher.Write([]byte(msg))

	if err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hasher.Sum(nil), sig); err != nil {
		return err
	}

	return nil
}
