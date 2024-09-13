package server

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type kvs struct {
	n   map[string]string
	ses map[string]string

	db *redis.Client
}

func initKVS() *kvs {
	return &kvs{
		n:   make(map[string]string),
		ses: make(map[string]string),
		db: redis.NewClient(&redis.Options{
			Addr: "redis.comame.dev:6379",
		}),
	}
}

func redisKey(name, key string) string {
	return fmt.Sprintf("note-comame-xyz:%s:%s", name, key)
}

func (k *kvs) SetNonce(key, nonce string, maxAge int) {
	k.db.Set(context.Background(), redisKey("nonce", key), nonce, time.Duration(maxAge)*time.Second)
}
func (k *kvs) GetNonce(key string) (string, error) {
	v, err := k.db.Get(context.Background(), redisKey("nonce", key)).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}
func (k *kvs) DeleteNonce(key string) {
	k.db.Del(context.Background(), redisKey("nonce", key))
}

func (k *kvs) SetSession(key, userID string) {
	k.db.Set(context.Background(), redisKey("session", key), userID, time.Duration(1*24)*time.Hour)
}
func (k *kvs) GetSession(key string) (userID string, ok bool) {
	v, err := k.db.Get(context.Background(), redisKey("session", key)).Result()
	if err != nil {
		return "", false
	}
	return v, true
}

func (k *kvs) DeleteSession(key string) {
	k.db.Del(context.Background(), redisKey("session", key))
}
