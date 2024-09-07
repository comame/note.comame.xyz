package server

import (
	"errors"
)

type kvs struct {
	n   map[string]string
	ses map[string]string
}

func initKVS() *kvs {
	return &kvs{
		n:   make(map[string]string),
		ses: make(map[string]string),
	}
}

func (k *kvs) SetNonce(key, nonce string, maxAge int) {
	k.n[key] = nonce
}
func (k *kvs) GetNonce(key string) (string, error) {
	v, ok := k.n[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}
func (k *kvs) DeleteNonce(key string) {
	delete(k.n, key)
}

func (k *kvs) SetSession(key, userID string) {
	k.ses[key] = userID
}
func (k *kvs) GetSession(key string) (userID string, ok bool) {
	u, ok := k.ses[key]
	if !ok {
		return "", false
	}
	return u, true
}
func (k *kvs) DeleteSession(key string) {
	delete(k.ses, key)
}
