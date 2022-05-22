package main

import (
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

var localCache CacheInterface

type AppCache struct {
	client *cache.Cache
}

type CacheBody struct {
	UserID    string        `json:"userId"`
	Limit     time.Duration `json:"limit"`
	Requests  uint64        `json:"requests"`
	ExpiresAt time.Time     `json:"expiresAt"`
}

type CacheInterface interface {
	Set(key string, data interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}

func (a AppCache) Set(key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	a.client.Set(key, b, expiration)
	return nil
}

func (a AppCache) Get(key string) ([]byte, error) {
	res, exist := a.client.Get(key)
	if !exist {
		return nil, nil
	}

	resByte, ok := res.([]byte)
	if !ok {
		return nil, errors.New("Format is not arr of bytes")
	}

	return resByte, nil
}

func InitMemCache() CacheInterface {
	localCache = &AppCache{
		// Create a cache with a default expiration time of 60 minutes, and which
		// purges expired items every 120 minutes
		client: cache.New(60*time.Minute, 120*time.Minute),
	}
	return localCache
}
