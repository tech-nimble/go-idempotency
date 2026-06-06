// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

package idempotency

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestGenerateHashKey(t *testing.T) {
	t.Parallel()

	t.Run("without header returns empty", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("POST", "/orders", nil)
		if got := generateHashKey(req); got != "" {
			t.Fatalf("expected empty key, got %q", got)
		}
	})

	t.Run("with header returns composite key", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("POST", "/orders", nil)
		req.Header.Set(headerIdempotencyKey, "abc")

		want := "abc_POST_/orders"
		if got := generateHashKey(req); got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	})
}

func TestNewIdempotencyOptions(t *testing.T) {
	t.Parallel()

	m := NewIdempotency(WithHeaderKey("X-Custom-Key"))
	if m.headerKey != "X-Custom-Key" {
		t.Fatalf("expected custom header key, got %q", m.headerKey)
	}

	// Each call must return an independent instance.
	other := NewIdempotency()
	if other.headerKey == "X-Custom-Key" {
		t.Fatal("instances share state")
	}
}

func TestEmptyStorage(t *testing.T) {
	t.Parallel()

	s := NewEmptyStorage()

	if _, err := s.Get(context.Background(), "k"); err == nil {
		t.Fatal("expected error from empty storage Get")
	}

	if err := s.Set(context.Background(), "k", []byte("v"), 0); err != nil {
		t.Fatalf("expected nil from empty storage Set, got %v", err)
	}
}
