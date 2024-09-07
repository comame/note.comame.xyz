package oidc

type kvs interface {
	SetNonce(key, nonce string, maxAge int)
	GetNonce(key string) (string, error)
	DeleteNonce(key string)
}
