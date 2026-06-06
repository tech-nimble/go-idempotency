// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package idempotency

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
)

// RedisStorage is a storage based on Redis.
type RedisStorage struct {
	redis *cache.Cache
}

// NewRedisStorage creates a new instance of RedisStorage.
func NewRedisStorage(redis *cache.Cache) *RedisStorage {
	return &RedisStorage{
		redis: redis,
	}
}

// Get returns data from Redis.
func (i *RedisStorage) Get(ctx context.Context, key string) ([]byte, error) {
	var result []byte

	err := i.redis.Get(ctx, key, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Set saves data to Redis.
func (i *RedisStorage) Set(ctx context.Context, key string, response []byte, expiration time.Duration) error {
	return i.redis.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: response,
		TTL:   expiration,
	})
}
