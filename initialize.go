// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package idempotency

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// RequestUniquenessChecker guards an endpoint against duplicate requests.
type RequestUniquenessChecker interface {
	Api(ctx *gin.Context)
}

// Initialize builds a Redis-backed idempotency middleware.
func Initialize(rl *redis.Client) RequestUniquenessChecker {
	storage := NewRedisStorage(
		cache.New(&cache.Options{
			Redis: rl,
		}),
	)

	return NewIdempotency(WithStorage(storage), WithHeaderKey(headerIdempotencyKey))
}
