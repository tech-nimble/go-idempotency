// SPDX-FileCopyrightText: 2025 Nimble Tech
// SPDX-License-Identifier: MIT

// Package idempotency provides a Gin middleware that makes endpoints idempotent
// by caching responses keyed by an idempotency header.
package idempotency

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	headerIdempotencyKey   = "X-Idempotency-Key"
	headerIdempotencyCache = "X-Idempotency-Cache"
	cacheHit               = "HIT"
)

const expiration = time.Hour

// Storage is an interface for working with storage.
type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, response []byte, expiration time.Duration) error
}

type bufferedWriter struct {
	gin.ResponseWriter
	Buffer bytes.Buffer
}

func (g *bufferedWriter) Write(data []byte) (int, error) {
	g.Buffer.Write(data)

	return g.ResponseWriter.Write(data)
}

type response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Options is a type for middleware options.
type Options func(*Idempotency)

// Idempotency is a middleware that checks the uniqueness of the request by the header key.
type Idempotency struct {
	storage   Storage
	headerKey string
}

// NewIdempotency creates a new instance of Idempotency. Without options it uses
// an empty (no-op) storage and the default idempotency header.
func NewIdempotency(options ...Options) *Idempotency {
	m := &Idempotency{
		storage:   NewEmptyStorage(),
		headerKey: headerIdempotencyKey,
	}

	for _, o := range options {
		o(m)
	}

	return m
}

// Api is a middleware that checks the uniqueness of the request by the header key.
//
//nolint:revive // method name kept for backward compatibility
func (i *Idempotency) Api(ctx *gin.Context) {
	key := generateHashKey(ctx.Request)
	if key == "" {
		err := errors.New("укажите в заголовке " + i.headerKey + " уникальный номер запроса в формате uuid v4 ")
		log.Error().Err(err)

		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			getErrorResponse(err, http.StatusBadRequest),
		)

		return
	}

	var resp *response
	respBytes, err := i.storage.Get(ctx, key)
	if err == nil {
		err = json.Unmarshal(respBytes, &resp)
		if err != nil {
			log.Error().Err(err).Interface("resp", resp).Msg("Could not marshal response to json.")

			ctx.Next()
			return
		}

		ctx.Header(headerIdempotencyCache, cacheHit)
		for k, v := range resp.Headers {
			ctx.Writer.Header()[k] = v
		}

		ctx.Status(resp.StatusCode)
		_, err = ctx.Writer.Write(resp.Body)
		if err != nil {
			log.Error().Err(err).Interface("resp", resp).Msg("Could not write response.")

			ctx.Next()
			return
		}

		ctx.Abort()
		return
	}

	buff := bytes.Buffer{}
	newWriter := &bufferedWriter{ctx.Writer, buff}

	ctx.Writer = newWriter
	ctx.Next()

	if !ctx.IsAborted() {
		return
	}

	if ctx.Writer.Status() >= http.StatusInternalServerError {
		return
	}

	resp, err = i.saveResponse(ctx, key, newWriter)
	if err != nil {
		log.Error().Err(err).Interface("resp", resp).Msg("Could not save response.")

		return
	}
}

// WithStorage sets the storage for the middleware.
func WithStorage(storage Storage) Options {
	return func(i *Idempotency) {
		i.storage = storage
	}
}

// WithHeaderKey sets the header key for the middleware.
func WithHeaderKey(headerKey string) Options {
	return func(i *Idempotency) {
		i.headerKey = headerKey
	}
}

func generateHashKey(r *http.Request) string {
	idempotencyKey := r.Header.Get(headerIdempotencyKey)
	if idempotencyKey == "" {
		return ""
	}

	return fmt.Sprintf("%s_%s_%s", idempotencyKey, r.Method, r.RequestURI)
}

func (i *Idempotency) saveResponse(ctx *gin.Context, hashKey string, writer *bufferedWriter) (*response, error) {
	resp := &response{
		StatusCode: ctx.Writer.Status(),
		Headers:    ctx.Writer.Header(),
		Body:       writer.Buffer.Bytes(),
	}

	content, err := json.Marshal(resp)
	if err != nil {
		return nil, errors.New("could not marshal response to json")
	}

	err = i.storage.Set(ctx, hashKey, content, expiration)
	if err != nil {
		return nil, errors.New("could not cache response")
	}

	return resp, nil
}

func getErrorResponse(err error, httpStatus int) jsonapi.ErrorsPayload {
	statusStr := strconv.Itoa(httpStatus)

	return jsonapi.ErrorsPayload{
		Errors: []*jsonapi.ErrorObject{
			{
				ID:     uuid.New().String(),
				Title:  err.Error(),
				Detail: err.Error(),
				Status: statusStr,
				Code:   statusStr,
			},
		},
	}
}
