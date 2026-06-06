// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package idempotency

import (
	"context"
	"errors"
	"time"
)

// EmptyStorage is an empty storage.
type EmptyStorage struct{}

// NewEmptyStorage creates a new instance of EmptyStorage.
func NewEmptyStorage() *EmptyStorage {
	return &EmptyStorage{}
}

// Get returns an error.
func (i *EmptyStorage) Get(_ context.Context, _ string) ([]byte, error) {
	return nil, errors.New("no data")
}

// Set returns nil.
func (i *EmptyStorage) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error {
	return nil
}
