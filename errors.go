package goredis

import "errors"

var (
	// ErrInvalidConfig indicates configuration validation failed.
	ErrInvalidConfig = errors.New("redis: invalid config")
	// ErrNotInitialized indicates a client was not created.
	ErrNotInitialized = errors.New("redis: not initialized")
)
