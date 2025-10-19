package cache

import (
	"context"
	"time"
)

// noOpCache is a cache implementation that does nothing.
// Useful for testing or when caching is disabled.
type noOpCache[T any] struct{}

// NewNoOp creates a new no-op cache.
// This is a convenience function for creating no-op caches directly.
func NewNoOp[T any]() Cache[T] {
	return &noOpCache[T]{}
}

func (c *noOpCache[T]) Get(
	_ context.Context,
	_ string,
) (T, bool) {
	var zero T
	return zero, false
}

func (c *noOpCache[T]) Set(
	_ context.Context,
	_ string,
	_ T,
	_ time.Duration,
) error {
	return nil
}

func (c *noOpCache[T]) Delete(
	_ context.Context,
	_ string,
) error {
	return nil
}

func (c *noOpCache[T]) Close() error {
	return nil
}
